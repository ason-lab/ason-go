package asun

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

// ---------------------------------------------------------------------------
// Decode — deserialize from ASUN
// ---------------------------------------------------------------------------

// Decode parses ASUN data and stores the result in v.
// v must be a pointer to a struct or a pointer to a slice of structs.
// Single struct input:  {field1,field2,...}:(val1,val2,...)
// Slice input:         [{field1,field2,...}]:(val1,val2,...),(val3,val4,...)
// The [{...}]: bracket wrapper is required for slice targets.
func Decode(data []byte, v any) error {
	d := decoder{
		data: data,
		pos:  0,
	}
	d.skipWhitespaceAndComments()

	// Detect target type and input format
	rv := reflect.ValueOf(v)
	isSliceTarget := rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Slice
	isAnyTarget := rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Interface
	startsWithBracket := d.pos < len(d.data) && d.data[d.pos] == '[' &&
		d.pos+1 < len(d.data) && d.data[d.pos+1] == '{'

	// Untyped target: bare value / plain array. SPEC §8.3 forbids top-level
	// `(...)` without a schema header.
	if isAnyTarget && !startsWithBracket {
		if d.pos < len(d.data) && d.data[d.pos] == '{' {
			return &UnmarshalError{d.pos, "schema header requires a typed target"}
		}
		if d.pos < len(d.data) && d.data[d.pos] == '(' {
			// `()` is the untyped null marker; longer bare tuples remain
			// an error per SPEC §8.3.
			if d.pos+1 < len(d.data) && d.data[d.pos+1] == ')' {
				d.pos += 2
				d.skipWhitespaceAndComments()
				if d.pos < len(d.data) {
					return &UnmarshalError{d.pos, "trailing characters after decoded value"}
				}
				// Leave the interface as zero-value (nil).
				return nil
			}
			return &UnmarshalError{d.pos, "bare tuple at top level — schema required"}
		}
		val, err := d.parseAnyValue()
		if err != nil {
			return err
		}
		d.skipWhitespaceAndComments()
		if d.pos < len(d.data) {
			return &UnmarshalError{d.pos, "trailing characters after decoded value"}
		}
		rv.Elem().Set(reflect.ValueOf(val))
		return nil
	}

	// Strict format enforcement: slices require [{...}]:, structs require {...}:
	if isSliceTarget && !startsWithBracket {
		return &UnmarshalError{d.pos, "Decode requires '[{...}]:' format for slice types"}
	}
	if startsWithBracket {
		return d.decodeSliceTop(v)
	}
	if err := d.unmarshalTop(v); err != nil {
		return err
	}
	d.skipWhitespaceAndComments()
	if d.pos < len(d.data) {
		return &UnmarshalError{d.pos, "trailing characters after decoded value"}
	}
	return nil
}

// decodeSliceTop parses [{schema}]:(v1,...),(v2,...) format
func (d *decoder) decodeSliceTop(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return &UnmarshalError{d.pos, "Decode requires a pointer to slice for [{...}]: format"}
	}
	sliceVal := rv.Elem()
	elemType := sliceVal.Type().Elem()
	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return &UnmarshalError{d.pos, "Decode requires a slice of structs"}
	}

	// Skip '['
	d.pos++

	// Parse schema
	if d.pos >= len(d.data) || d.data[d.pos] != '{' {
		return d.errorf("expected '{'")
	}
	fields, schemaKey, err := d.parseSchema()
	if err != nil {
		return err
	}
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.data[d.pos] != ']' {
		return d.errorf("expected ']'")
	}
	d.pos++
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.data[d.pos] != ':' {
		return d.errorf("expected ':'")
	}
	d.pos++

	si := getStructInfo(elemType)
	fieldMap := buildFieldMapCached(si, fields, schemaKey)
	if sliceVal.Type().Elem().Kind() == reflect.Struct {
		if rowCount := countTopLevelTupleRows(d.data, d.pos); rowCount > 0 {
			if sliceVal.Cap() < rowCount {
				sliceVal = reflect.MakeSlice(sliceVal.Type(), rowCount, rowCount)
			} else {
				sliceVal = sliceVal.Slice(0, rowCount)
			}
			exact := isIdentityFieldMap(si, fieldMap)
			decoded := 0
			for decoded < rowCount {
				d.skipWhitespaceAndComments()
				if d.pos >= len(d.data) || d.data[d.pos] != '(' {
					break
				}
				elem := sliceVal.Index(decoded)
				if exact {
					err = d.unmarshalTupleExact(elem, si)
				} else {
					err = d.unmarshalTuple(elem, si, fieldMap)
				}
				if err != nil {
					return err
				}
				decoded++
				d.skipWhitespaceAndComments()
				if d.pos < len(d.data) && d.data[d.pos] == ',' {
					d.pos++
					d.skipWhitespaceAndComments()
					if d.pos >= len(d.data) || d.data[d.pos] != '(' {
						break
					}
				}
			}
			rv.Elem().Set(sliceVal.Slice(0, decoded))
			return nil
		}
	}

	// Parse rows
	for {
		d.skipWhitespaceAndComments()
		if d.pos >= len(d.data) {
			break
		}
		if d.data[d.pos] != '(' {
			break
		}

		elem := reflect.New(elemType).Elem()
		if err := d.unmarshalTuple(elem, si, fieldMap); err != nil {
			return err
		}
		sliceVal = reflect.Append(sliceVal, elem)

		d.skipWhitespaceAndComments()
		if d.pos < len(d.data) && d.data[d.pos] == ',' {
			d.pos++
			d.skipWhitespaceAndComments()
			if d.pos >= len(d.data) || d.data[d.pos] != '(' {
				break
			}
		}
	}

	rv.Elem().Set(sliceVal)
	return nil
}

// ---------------------------------------------------------------------------
// decoder
// ---------------------------------------------------------------------------

type decoder struct {
	data  []byte
	pos   int
	depth int
}

// maxDepth bounds recursion while parsing nested structures (schema
// annotations, arrays, tuples, nested structs). Untrusted input could
// otherwise drive unbounded recursion into a fatal, unrecoverable stack
// overflow that aborts the whole process (a DoS). 128 is well beyond any
// realistic nesting yet far below the point where Go's growable goroutine
// stack is at risk.
const maxDepth = 128

// enter increments the recursion depth, returning an error if the limit is
// exceeded. Each successful enter must be paired with a leave (use defer).
func (d *decoder) enter() error {
	d.depth++
	if d.depth > maxDepth {
		return d.errorf("maximum nesting depth exceeded")
	}
	return nil
}

func (d *decoder) leave() {
	d.depth--
}

// boundedCache is a concurrency-safe schema cache with a hard entry cap.
// Without a cap, every distinct schema in (untrusted) input is retained
// forever — an unbounded memory-growth DoS. When the cap is reached the whole
// map is dropped and repopulated; legitimate workloads use a small, stable set
// of schemas and never hit the cap, while adversarial input degrades only to
// occasional cache misses (correctness is unaffected). The hot path stays a
// lock-free sync.Map load.
type boundedCache struct {
	m   sync.Map
	n   atomic.Int64
	max int64
}

func (c *boundedCache) Load(key string) (any, bool) {
	return c.m.Load(key)
}

// LoadOrStore returns the existing value for key if present; otherwise stores
// and returns value. It flushes the cache before inserting if the cap is hit.
func (c *boundedCache) LoadOrStore(key string, value any) (any, bool) {
	if actual, ok := c.m.Load(key); ok {
		return actual, true
	}
	max := c.max
	if max <= 0 {
		max = maxCachedSchemas
	}
	if c.n.Load() >= max {
		c.m.Range(func(k, _ any) bool {
			c.m.Delete(k)
			c.n.Add(-1)
			return true
		})
	}
	actual, loaded := c.m.LoadOrStore(key, value)
	if !loaded {
		c.n.Add(1)
	}
	return actual, loaded
}

func (c *boundedCache) Store(key string, value any) {
	c.LoadOrStore(key, value)
}

// maxCachedSchemas caps distinct schema headers retained per cache.
const maxCachedSchemas = 4096

var schemaFieldsCache = boundedCache{max: maxCachedSchemas}

func countTopLevelEntries(data []byte, start int, open, close byte) int {
	if start >= len(data) || data[start] != open {
		return 0
	}
	depthParen, depthBracket := 0, 0
	inString := false
	count := 0
	hasValue := false
	for i := start + 1; i < len(data); i++ {
		b := data[i]
		if inString {
			if b == '\\' && i+1 < len(data) {
				i++
				continue
			}
			if b == '"' {
				inString = false
			}
			hasValue = true
			continue
		}
		if b == '"' {
			inString = true
			hasValue = true
			continue
		}
		if b == '/' && i+1 < len(data) && data[i+1] == '*' {
			i += 2
			for i+1 < len(data) && !(data[i] == '*' && data[i+1] == '/') {
				i++
			}
			if i+1 < len(data) {
				i++
			}
			continue
		}
		switch b {
		case '(':
			depthParen++
			hasValue = true
		case ')':
			if depthParen > 0 {
				depthParen--
			}
		case '[':
			depthBracket++
			hasValue = true
		case ']':
			if depthBracket > 0 {
				depthBracket--
			} else if close == ']' {
				if hasValue {
					count++
				}
				return count
			}
		case ',':
			if depthParen == 0 && depthBracket == 0 {
				if hasValue {
					count++
					hasValue = false
				}
			} else {
				hasValue = true
			}
		default:
			if b == close && depthParen == 0 && depthBracket == 0 {
				if hasValue {
					count++
				}
				return count
			}
			switch b {
			case ' ', '\t', '\n', '\r':
			default:
				hasValue = true
			}
		}
	}
	return count
}

func countTopLevelTupleRows(data []byte, start int) int {
	depthBracket := 0
	inString := false
	count := 0
	for i := start; i < len(data); i++ {
		b := data[i]
		if inString {
			if b == '\\' && i+1 < len(data) {
				i++
				continue
			}
			if b == '"' {
				inString = false
			}
			continue
		}
		if b == '"' {
			inString = true
			continue
		}
		if b == '/' && i+1 < len(data) && data[i+1] == '*' {
			i += 2
			for i+1 < len(data) && !(data[i] == '*' && data[i+1] == '/') {
				i++
			}
			if i+1 < len(data) {
				i++
			}
			continue
		}
		switch b {
		case '[':
			depthBracket++
		case ']':
			if depthBracket > 0 {
				depthBracket--
			}
		case '(':
			if depthBracket == 0 {
				count++
			}
		}
	}
	return count
}

func (d *decoder) errorf(format string, args ...any) *UnmarshalError {
	msg := format
	if len(args) > 0 {
		msg = sprintf(format, args...)
	}
	return &UnmarshalError{Pos: d.pos, Message: msg}
}

func sprintf(format string, args ...any) string {
	// Minimal sprintf to avoid fmt dependency in hot path
	result := format
	for _, a := range args {
		switch v := a.(type) {
		case string:
			result = strings.Replace(result, "%s", v, 1)
		case byte:
			result = strings.Replace(result, "%c", string(rune(v)), 1)
		}
	}
	return result
}

func (d *decoder) skipWhitespace() {
	for d.pos < len(d.data) {
		switch d.data[d.pos] {
		case ' ', '\t', '\n', '\r':
			d.pos++
		default:
			return
		}
	}
}

func (d *decoder) skipWhitespaceAndComments() {
	for {
		d.skipWhitespace()
		if d.pos+1 < len(d.data) && d.data[d.pos] == '/' && d.data[d.pos+1] == '*' {
			d.pos += 2
			for d.pos+1 < len(d.data) {
				if d.data[d.pos] == '*' && d.data[d.pos+1] == '/' {
					d.pos += 2
					break
				}
				d.pos++
			}
		} else {
			return
		}
	}
}

// parseSchema parses {field1,field2,...} or {field1@type,field2@{...},field3@[...],...}
// Returns field names only (type annotations are skipped).
func skipInlineWhitespace(data []byte, pos int) int {
	for pos < len(data) && (data[pos] == ' ' || data[pos] == '\t') {
		pos++
	}
	return pos
}

func scanBalancedEnd(data []byte, pos int, open, close byte) (int, error) {
	depth := 0
	for pos < len(data) {
		b := data[pos]
		pos++
		if b == open {
			depth++
		} else if b == close {
			depth--
			if depth == 0 {
				return pos, nil
			}
		}
	}
	return 0, &UnmarshalError{pos, "unbalanced brackets"}
}

func scanSchemaEnd(data []byte, pos int) (int, error) {
	if pos >= len(data) || data[pos] != '{' {
		return 0, &UnmarshalError{pos, "expected '{'"}
	}
	pos++
	fieldCount := 0
	for {
		pos = skipInlineWhitespace(data, pos)
		if pos >= len(data) {
			return 0, &UnmarshalError{pos, "unexpected EOF in schema"}
		}
		if data[pos] == '}' {
			return pos + 1, nil
		}
		if fieldCount > 0 {
			if data[pos] != ',' {
				return 0, &UnmarshalError{pos, "expected ','"}
			}
			pos++
			pos = skipInlineWhitespace(data, pos)
		}
		if pos < len(data) && data[pos] == '"' {
			pos++
			for pos < len(data) {
				b := data[pos]
				if b == '\\' && pos+1 < len(data) {
					pos += 2
					continue
				}
				pos++
				if b == '"' {
					break
				}
			}
		} else {
			for pos < len(data) {
				b := data[pos]
				if b == ',' || b == '}' || b == '@' || b == ':' || b == ' ' || b == '\t' {
					break
				}
				pos++
			}
		}
		pos = skipInlineWhitespace(data, pos)
		if pos < len(data) && data[pos] == '@' {
			pos++
			pos = skipInlineWhitespace(data, pos)
			if pos >= len(data) {
				return 0, &UnmarshalError{pos, "unexpected EOF in schema"}
			}
			switch data[pos] {
			case '{':
				end, err := scanBalancedEnd(data, pos, '{', '}')
				if err != nil {
					return 0, err
				}
				pos = end
			case '[':
				end, err := scanBalancedEnd(data, pos, '[', ']')
				if err != nil {
					return 0, err
				}
				pos = end
			default:
				for pos < len(data) {
					b := data[pos]
					if b == ',' || b == '}' || b == ' ' || b == '\t' {
						break
					}
					pos++
				}
			}
		}
		fieldCount++
	}
}

func (d *decoder) parseSchema() ([]string, string, error) {
	if d.data[d.pos] != '{' {
		return nil, "", d.errorf("expected '{'")
	}
	start := d.pos
	if end, err := scanSchemaEnd(d.data, start); err == nil {
		key := unsafeString(d.data[start:end])
		if cached, ok := schemaFieldsCache.Load(key); ok {
			d.pos = end
			return cached.([]string), key, nil
		}
	}
	d.pos++

	var fields []string
	for {
		d.skipWhitespace()
		if d.pos >= len(d.data) {
			return nil, "", d.errorf("unexpected EOF in schema")
		}
		if d.data[d.pos] == '}' {
			d.pos++
			break
		}
		if len(fields) > 0 {
			if d.data[d.pos] != ',' {
				return nil, "", d.errorf("expected ','")
			}
			d.pos++
			d.skipWhitespace()
		}

		// Parse field name
		var name string
		if d.pos < len(d.data) && d.data[d.pos] == '"' {
			parsed, err := d.parseQuotedString()
			if err != nil {
				return nil, "", err
			}
			name = parsed
		} else {
			start := d.pos
			for d.pos < len(d.data) {
				b := d.data[d.pos]
				if b == ',' || b == '}' || b == '@' || b == ':' || b == ' ' || b == '\t' {
					break
				}
				d.pos++
			}
			name = unsafeString(d.data[start:d.pos])
		}
		d.skipWhitespace()

		// Validate and skip optional type annotation after '@'
		if d.pos < len(d.data) && d.data[d.pos] == '@' {
			d.pos++
			d.skipWhitespace()
			if err := d.parseSchemaAnnotation(); err != nil {
				return nil, "", err
			}
		}

		fields = append(fields, name)
	}
	key := unsafeString(d.data[start:d.pos])
	schemaFieldsCache.Store(key, fields)
	return fields, key, nil
}

func (d *decoder) parseSchemaAnnotation() error {
	if err := d.enter(); err != nil {
		return err
	}
	defer d.leave()
	if d.pos >= len(d.data) {
		return d.errorf("expected schema type after '@'")
	}
	switch d.data[d.pos] {
	case '{':
		_, _, err := d.parseSchema()
		return err
	case '[':
		d.pos++
		d.skipWhitespace()
		if d.pos < len(d.data) && d.data[d.pos] == ']' {
			d.pos++
			return nil
		}
		if d.pos < len(d.data) && d.data[d.pos] == '{' {
			if _, _, err := d.parseSchema(); err != nil {
				return err
			}
		} else {
			if err := d.parseAllowedSchemaScalarType(); err != nil {
				return err
			}
		}
		d.skipWhitespace()
		if d.pos >= len(d.data) || d.data[d.pos] != ']' {
			return d.errorf("expected ']' in array type annotation")
		}
		d.pos++
		return nil
	default:
		return d.parseAllowedSchemaScalarType()
	}
}

func (d *decoder) parseAllowedSchemaScalarType() error {
	start := d.pos
	for d.pos < len(d.data) {
		b := d.data[d.pos]
		if b == ',' || b == '}' || b == ']' || b == ' ' || b == '\t' {
			break
		}
		d.pos++
	}
	if start == d.pos {
		return d.errorf("expected schema type after '@'")
	}
	token := unsafeString(d.data[start:d.pos])
	if strings.HasSuffix(token, "?") {
		token = token[:len(token)-1]
	}
	switch token {
	case "int", "str", "float", "bool":
		return nil
	default:
		return d.errorf("unsupported schema type %q; use int, str, float, or bool", token)
	}
}

func (d *decoder) skipBalanced(open, close byte) error {
	depth := 0
	for d.pos < len(d.data) {
		b := d.data[d.pos]
		d.pos++
		if b == open {
			depth++
		} else if b == close {
			depth--
			if depth == 0 {
				return nil
			}
		}
	}
	return d.errorf("unbalanced brackets")
}

// ---------------------------------------------------------------------------
// buildFieldMap maps schema field names to struct field indices
// ---------------------------------------------------------------------------

func buildFieldMap(si *structInfo, schemaFields []string) []int {
	if len(schemaFields) == len(si.fields) {
		exact := true
		for i, name := range schemaFields {
			if si.fields[i].name != name {
				exact = false
				break
			}
		}
		if exact {
			return si.identityFieldMap
		}
	}
	m := make([]int, len(schemaFields))
	for i, name := range schemaFields {
		if idx, ok := si.nameIndex[name]; ok {
			m[i] = idx
		} else {
			m[i] = -1
		}
	}
	return m
}

func buildFieldMapCached(si *structInfo, schemaFields []string, schemaKey string) []int {
	if schemaKey != "" {
		if cached, ok := si.fieldMapCache.Load(schemaKey); ok {
			return cached.([]int)
		}
	}
	fieldMap := buildFieldMap(si, schemaFields)
	if schemaKey != "" {
		actual, _ := si.fieldMapCache.LoadOrStore(schemaKey, fieldMap)
		return actual.([]int)
	}
	return fieldMap
}

func isIdentityFieldMap(si *structInfo, fieldMap []int) bool {
	if len(fieldMap) != len(si.identityFieldMap) {
		return false
	}
	if len(fieldMap) == 0 {
		return true
	}
	return &fieldMap[0] == &si.identityFieldMap[0]
}

func fieldByInfo(rv reflect.Value, fi fieldInfo) reflect.Value {
	if fi.direct >= 0 {
		return rv.Field(fi.direct)
	}
	return rv.FieldByIndex(fi.index)
}

// ---------------------------------------------------------------------------
// unmarshalTop — parse schema header then data for a single struct
// ---------------------------------------------------------------------------

func (d *decoder) unmarshalTop(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return d.errorf("Unmarshal requires a pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return d.errorf("Unmarshal requires a pointer to struct")
	}

	si := getStructInfo(rv.Type())

	// Parse schema
	if d.pos >= len(d.data) || d.data[d.pos] != '{' {
		return d.errorf("expected '{'")
	}
	fields, schemaKey, err := d.parseSchema()
	if err != nil {
		return err
	}

	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.data[d.pos] != ':' {
		return d.errorf("expected ':'")
	}
	d.pos++
	d.skipWhitespaceAndComments()

	fieldMap := buildFieldMapCached(si, fields, schemaKey)
	return d.unmarshalTuple(rv, si, fieldMap)
}

// ---------------------------------------------------------------------------
// unmarshalTuple — parse (val1,val2,...) positionally
// ---------------------------------------------------------------------------

func (d *decoder) unmarshalTuple(rv reflect.Value, si *structInfo, fieldMap []int) error {
	if isIdentityFieldMap(si, fieldMap) {
		return d.unmarshalTupleExact(rv, si)
	}
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.data[d.pos] != '(' {
		return d.errorf("expected '('")
	}
	d.pos++

	for i := 0; i < len(fieldMap); i++ {
		d.skipWhitespaceAndComments()
		if d.pos >= len(d.data) {
			return d.errorf("unexpected EOF in tuple")
		}
		if d.data[d.pos] == ')' {
			break
		}
		if i > 0 {
			if d.data[d.pos] == ',' {
				d.pos++
				d.skipWhitespaceAndComments()
				if d.pos < len(d.data) && d.data[d.pos] == ')' {
					break
				}
			} else if d.data[d.pos] == ')' {
				break
			} else {
				return d.errorf("expected ',' or ')'")
			}
		}

		fi := fieldMap[i]
		if fi < 0 {
			// Unknown field — skip value
			if err := d.skipValue(); err != nil {
				return err
			}
			continue
		}

		fv := fieldByInfo(rv, si.fields[fi])
		if err := d.unmarshalValue(fv); err != nil {
			return err
		}
	}

	// Skip any remaining values (source struct may have more fields than target)
	d.skipWhitespaceAndComments()
	for d.pos < len(d.data) && d.data[d.pos] != ')' {
		if d.data[d.pos] == ',' {
			d.pos++
			d.skipWhitespaceAndComments()
			if d.pos < len(d.data) && d.data[d.pos] == ')' {
				break
			}
		}
		if d.pos < len(d.data) && d.data[d.pos] != ')' {
			if err := d.skipValue(); err != nil {
				return err
			}
			d.skipWhitespaceAndComments()
		}
	}

	if d.pos < len(d.data) && d.data[d.pos] == ')' {
		d.pos++
	}
	return nil
}

func (d *decoder) unmarshalTupleExact(rv reflect.Value, si *structInfo) error {
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.data[d.pos] != '(' {
		return d.errorf("expected '('")
	}
	d.pos++

	for i := 0; i < len(si.fields); i++ {
		d.skipWhitespaceAndComments()
		if d.pos >= len(d.data) {
			return d.errorf("unexpected EOF in tuple")
		}
		if d.data[d.pos] == ')' {
			break
		}
		if i > 0 {
			if d.data[d.pos] == ',' {
				d.pos++
				d.skipWhitespaceAndComments()
				if d.pos < len(d.data) && d.data[d.pos] == ')' {
					break
				}
			} else if d.data[d.pos] == ')' {
				break
			} else {
				return d.errorf("expected ',' or ')'")
			}
		}

		fv := fieldByInfo(rv, si.fields[i])
		if err := d.unmarshalValue(fv); err != nil {
			return err
		}
	}

	d.skipWhitespaceAndComments()
	for d.pos < len(d.data) && d.data[d.pos] != ')' {
		if d.data[d.pos] == ',' {
			d.pos++
			d.skipWhitespaceAndComments()
			if d.pos < len(d.data) && d.data[d.pos] == ')' {
				break
			}
		}
		if d.pos < len(d.data) && d.data[d.pos] != ')' {
			if err := d.skipValue(); err != nil {
				return err
			}
			d.skipWhitespaceAndComments()
		}
	}

	if d.pos < len(d.data) && d.data[d.pos] == ')' {
		d.pos++
	}
	return nil
}

// ---------------------------------------------------------------------------
// unmarshalValue — dispatch based on target field type
// ---------------------------------------------------------------------------

func (d *decoder) unmarshalValue(fv reflect.Value) error {
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) {
		return nil
	}

	// Handle pointer types
	if fv.Kind() == reflect.Ptr {
		if d.atValueEnd() {
			// nil
			fv.Set(reflect.Zero(fv.Type()))
			return nil
		}
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}
		return d.unmarshalValue(fv.Elem())
	}

	// Handle interface{}
	if fv.Kind() == reflect.Interface {
		val, err := d.parseAnyValue()
		if err != nil {
			return err
		}
		if val != nil {
			fv.Set(reflect.ValueOf(val))
		}
		return nil
	}

	switch fv.Kind() {
	case reflect.Bool:
		v, err := d.parseBool()
		if err != nil {
			return err
		}
		fv.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := d.parseInt64()
		if err != nil {
			return err
		}
		// reflect.Value.SetInt panics if v does not fit the sized field
		// (e.g. 200 into an int8). Reject out-of-range values instead.
		if fv.OverflowInt(v) {
			return &UnmarshalError{d.pos, "integer overflows target field"}
		}
		fv.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := d.parseUint64()
		if err != nil {
			return err
		}
		if fv.OverflowUint(v) {
			return &UnmarshalError{d.pos, "integer overflows target field"}
		}
		fv.SetUint(v)

	case reflect.Float32, reflect.Float64:
		v, err := d.parseFloat64()
		if err != nil {
			return err
		}
		fv.SetFloat(v)

	case reflect.String:
		s, err := d.parseString()
		if err != nil {
			return err
		}
		fv.SetString(s)

	case reflect.Slice:
		return d.unmarshalSlice(fv)

	case reflect.Map:
		return d.errorf("map fields are not supported")

	case reflect.Struct:
		return d.unmarshalNestedStruct(fv)

	default:
		return d.errorf("unsupported type: %s", fv.Type().String())
	}
	return nil
}

// ---------------------------------------------------------------------------
// Parsing primitives
// ---------------------------------------------------------------------------

func (d *decoder) atValueEnd() bool {
	if d.pos >= len(d.data) {
		return true
	}
	b := d.data[d.pos]
	return b == ',' || b == ')' || b == ']'
}

func (d *decoder) parseBool() (bool, error) {
	d.skipWhitespaceAndComments()
	if d.pos+4 <= len(d.data) &&
		d.data[d.pos] == 't' &&
		d.data[d.pos+1] == 'r' &&
		d.data[d.pos+2] == 'u' &&
		d.data[d.pos+3] == 'e' {
		if d.pos+4 >= len(d.data) || isDelim(d.data[d.pos+4]) {
			d.pos += 4
			return true, nil
		}
	}
	if d.pos+5 <= len(d.data) &&
		d.data[d.pos] == 'f' &&
		d.data[d.pos+1] == 'a' &&
		d.data[d.pos+2] == 'l' &&
		d.data[d.pos+3] == 's' &&
		d.data[d.pos+4] == 'e' {
		if d.pos+5 >= len(d.data) || isDelim(d.data[d.pos+5]) {
			d.pos += 5
			return false, nil
		}
	}
	return false, d.errorf("invalid bool")
}

func isDelim(b byte) bool {
	return b == ',' || b == ')' || b == ']' || b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func (d *decoder) parseInt64() (int64, error) {
	d.skipWhitespaceAndComments()
	neg := false
	if d.pos < len(d.data) && d.data[d.pos] == '-' {
		neg = true
		d.pos++
	}
	var val uint64
	digits := 0
	for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
		digit := uint64(d.data[d.pos] - '0')
		// Detect overflow instead of silently wrapping (data corruption).
		if val > (^uint64(0)-digit)/10 {
			return 0, d.errorf("integer overflow")
		}
		val = val*10 + digit
		d.pos++
		digits++
	}
	if digits == 0 {
		return 0, d.errorf("invalid number")
	}
	if neg {
		// Magnitude fits in int64 (>= math.MinInt64) exactly when val <= 2^63.
		if val > 1<<63 {
			return 0, d.errorf("integer overflow")
		}
		return -int64(val), nil
	}
	if val > 1<<63-1 {
		return 0, d.errorf("integer overflow")
	}
	return int64(val), nil
}

func (d *decoder) parseUint64() (uint64, error) {
	d.skipWhitespaceAndComments()
	var val uint64
	digits := 0
	for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
		digit := uint64(d.data[d.pos] - '0')
		if val > (^uint64(0)-digit)/10 {
			return 0, d.errorf("integer overflow")
		}
		val = val*10 + digit
		d.pos++
		digits++
	}
	if digits == 0 {
		return 0, d.errorf("invalid number")
	}
	return val, nil
}

func (d *decoder) parseFloat64() (float64, error) {
	d.skipWhitespaceAndComments()
	start := d.pos
	neg := false
	if d.pos < len(d.data) && d.data[d.pos] == '-' {
		neg = true
		d.pos++
	}
	intStart := d.pos
	for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
		d.pos++
	}
	if d.pos == intStart && neg {
		return 0, d.errorf("invalid number")
	}
	hasDot := d.pos < len(d.data) && d.data[d.pos] == '.'
	if hasDot {
		d.pos++
		for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
			d.pos++
		}
	}
	if d.pos == start || (d.pos == start+1 && neg) {
		return 0, d.errorf("invalid number")
	}
	// Fast path: integer-valued (no fractional part, no exponent). Fractional
	// values must NOT be accumulated as frac += digit*0.1^k — that is not
	// correctly-rounded (e.g. "0.3" → 0.30000000000000004, "2.675" →
	// 2.6750000000000003) and breaks round-trip. Those go through
	// strconv.ParseFloat below, which produces the nearest double.
	hasExp := d.pos < len(d.data) && (d.data[d.pos] == 'e' || d.data[d.pos] == 'E')
	if !hasExp && !hasDot && d.pos-start <= 18 {
		var intPart uint64
		p := start
		if neg {
			p++
		}
		for p < d.pos {
			intPart = intPart*10 + uint64(d.data[p]-'0')
			p++
		}
		v := float64(intPart)
		if neg {
			v = -v
		}
		return v, nil
	}
	// Fallback: fractional values, exponents, or very long numbers.
	if hasExp {
		d.pos++
		if d.pos < len(d.data) && (d.data[d.pos] == '+' || d.data[d.pos] == '-') {
			d.pos++
		}
		for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
			d.pos++
		}
	}
	s := unsafeString(d.data[start:d.pos])
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, d.errorf("invalid float: %s", s)
	}
	return v, nil
}

func (d *decoder) parseString() (string, error) {
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) {
		return "", nil
	}
	if d.data[d.pos] == '"' {
		return d.parseQuotedString()
	}
	return d.parsePlainValue()
}

// parsePlainValue reads unquoted string until delimiter, returns zerocopy slice.
func (d *decoder) parsePlainValue() (string, error) {
	start := d.pos
	hasEscape := false
	for d.pos < len(d.data) {
		b := d.data[d.pos]
		if b == ',' || b == ')' || b == ']' {
			break
		}
		if b == '\\' {
			hasEscape = true
			d.pos += 2
			continue
		}
		d.pos++
	}
	raw := d.data[start:d.pos]
	// Trim trailing whitespace
	for len(raw) > 0 && (raw[len(raw)-1] == ' ' || raw[len(raw)-1] == '\t') {
		raw = raw[:len(raw)-1]
	}
	// Trim leading whitespace
	for len(raw) > 0 && (raw[0] == ' ' || raw[0] == '\t') {
		raw = raw[1:]
	}
	if hasEscape {
		return unescapePlain(raw)
	}
	// Zerocopy
	return unsafeString(raw), nil
}

// parseQuotedString handles "..." with escape sequences.
// Zerocopy when no escapes found.
func (d *decoder) parseQuotedString() (string, error) {
	d.pos++ // skip opening "
	start := d.pos

	// Fast scan for closing quote without escapes
	scan := d.pos
	for scan < len(d.data) {
		if d.data[scan] == '"' {
			// No escapes — zerocopy
			s := unsafeString(d.data[start:scan])
			d.pos = scan + 1
			return s, nil
		}
		if d.data[scan] == '\\' {
			break
		}
		scan++
	}

	// Slow path: has escapes
	var buf []byte
	if scan > start {
		buf = append(buf, d.data[start:scan]...)
	}
	d.pos = scan

	for d.pos < len(d.data) {
		b := d.data[d.pos]
		if b == '"' {
			d.pos++
			return string(buf), nil
		}
		if b == '\\' {
			d.pos++
			if d.pos >= len(d.data) {
				return "", d.errorf("unclosed string")
			}
			esc := d.data[d.pos]
			d.pos++
			switch esc {
			case '"':
				buf = append(buf, '"')
			case '\\':
				buf = append(buf, '\\')
			case 'n':
				buf = append(buf, '\n')
			case 't':
				buf = append(buf, '\t')
			case 'r':
				buf = append(buf, '\r')
			case 'b':
				buf = append(buf, '\b')
			case 'f':
				buf = append(buf, '\f')
			case ',':
				buf = append(buf, ',')
			case '(':
				buf = append(buf, '(')
			case ')':
				buf = append(buf, ')')
			case '[':
				buf = append(buf, '[')
			case ']':
				buf = append(buf, ']')
			case '{':
				buf = append(buf, '{')
			case '}':
				buf = append(buf, '}')
			case '<':
				buf = append(buf, '<')
			case '>':
				buf = append(buf, '>')
			case ':':
				buf = append(buf, ':')
			case '@':
				buf = append(buf, '@')
			case 'u':
				if d.pos+4 > len(d.data) {
					return "", d.errorf("invalid unicode escape")
				}
				hexStr := unsafeString(d.data[d.pos : d.pos+4])
				cp, err := strconv.ParseUint(hexStr, 16, 32)
				if err != nil {
					return "", d.errorf("invalid unicode escape")
				}
				d.pos += 4
				r := rune(cp)
				// Combine a UTF-16 surrogate pair (\uD800-\uDBFF followed by
				// \uDC00-\uDFFF) into the single astral code point it encodes.
				if cp >= 0xD800 && cp <= 0xDBFF {
					// A high surrogate MUST be followed by a \uXXXX low
					// surrogate; a lone or mispaired high surrogate is invalid
					// and must be rejected (not silently turned into U+FFFD).
					if d.pos+6 > len(d.data) || d.data[d.pos] != '\\' || d.data[d.pos+1] != 'u' {
						return "", d.errorf("invalid unicode escape: unpaired surrogate")
					}
					lo, err := strconv.ParseUint(unsafeString(d.data[d.pos+2:d.pos+6]), 16, 32)
					if err != nil || lo < 0xDC00 || lo > 0xDFFF {
						return "", d.errorf("invalid unicode escape: unpaired surrogate")
					}
					r = ((rune(cp) - 0xD800) << 10) + (rune(lo) - 0xDC00) + 0x10000
					d.pos += 6
				} else if cp >= 0xDC00 && cp <= 0xDFFF {
					// A lone low surrogate is invalid.
					return "", d.errorf("invalid unicode escape: unpaired surrogate")
				}
				buf = append(buf, string(r)...)
			default:
				return "", d.errorf("invalid escape: \\%c", esc)
			}
		} else {
			buf = append(buf, b)
			d.pos++
		}
	}
	return "", d.errorf("unclosed string")
}

func unescapePlain(raw []byte) (string, error) {
	buf := make([]byte, 0, len(raw))
	i := 0
	for i < len(raw) {
		if raw[i] == '\\' {
			i++
			if i >= len(raw) {
				return "", &UnmarshalError{Message: "unexpected EOF in escape"}
			}
			switch raw[i] {
			case ',':
				buf = append(buf, ',')
			case '(':
				buf = append(buf, '(')
			case ')':
				buf = append(buf, ')')
			case '[':
				buf = append(buf, '[')
			case ']':
				buf = append(buf, ']')
			case '<':
				buf = append(buf, '<')
			case '>':
				buf = append(buf, '>')
			case ':':
				buf = append(buf, ':')
			case '"':
				buf = append(buf, '"')
			case '\\':
				buf = append(buf, '\\')
			case 'n':
				buf = append(buf, '\n')
			case 't':
				buf = append(buf, '\t')
			case 'r':
				buf = append(buf, '\r')
			case 'b':
				buf = append(buf, '\b')
			case 'f':
				buf = append(buf, '\f')
			case '{':
				buf = append(buf, '{')
			case '}':
				buf = append(buf, '}')
			case '@':
				buf = append(buf, '@')
			case 'u':
				if i+5 > len(raw) {
					return "", &UnmarshalError{Message: "invalid unicode escape"}
				}
				hexStr := unsafeString(raw[i+1 : i+5])
				cp, err := strconv.ParseUint(hexStr, 16, 32)
				if err != nil {
					return "", &UnmarshalError{Message: "invalid unicode escape"}
				}
				i += 4
				r := rune(cp)
				// Combine a UTF-16 surrogate pair into its astral code point.
				if cp >= 0xD800 && cp <= 0xDBFF {
					// A high surrogate MUST be followed by a \uXXXX low
					// surrogate; a lone or mispaired high surrogate is invalid.
					if i+7 > len(raw) || raw[i+1] != '\\' || raw[i+2] != 'u' {
						return "", &UnmarshalError{Message: "invalid unicode escape: unpaired surrogate"}
					}
					lo, err := strconv.ParseUint(unsafeString(raw[i+3:i+7]), 16, 32)
					if err != nil || lo < 0xDC00 || lo > 0xDFFF {
						return "", &UnmarshalError{Message: "invalid unicode escape: unpaired surrogate"}
					}
					r = ((rune(cp) - 0xD800) << 10) + (rune(lo) - 0xDC00) + 0x10000
					i += 6
				} else if cp >= 0xDC00 && cp <= 0xDFFF {
					// A lone low surrogate is invalid.
					return "", &UnmarshalError{Message: "invalid unicode escape: unpaired surrogate"}
				}
				buf = append(buf, string(r)...)
			default:
				return "", &UnmarshalError{Message: "invalid escape: \\" + string(rune(raw[i]))}
			}
		} else {
			buf = append(buf, raw[i])
		}
		i++
	}
	return string(buf), nil
}

// ---------------------------------------------------------------------------
// parseAnyValue — parse value when target type is interface{}
// ---------------------------------------------------------------------------

func (d *decoder) parseAnyValue() (any, error) {
	if err := d.enter(); err != nil {
		return nil, err
	}
	defer d.leave()
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.atValueEnd() {
		return nil, nil
	}
	b := d.data[d.pos]

	switch {
	case b == '"':
		return d.parseQuotedString()
	case b == '[':
		var arr []any
		d.pos++
		first := true
		for {
			d.skipWhitespaceAndComments()
			if d.pos >= len(d.data) || d.data[d.pos] == ']' {
				d.pos++
				break
			}
			if !first {
				if d.data[d.pos] == ',' {
					d.pos++
					d.skipWhitespaceAndComments()
					if d.pos < len(d.data) && d.data[d.pos] == ']' {
						d.pos++
						break
					}
				}
			}
			first = false
			v, err := d.parseAnyValue()
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	case b == '(':
		// `()` is the untyped null marker.
		if d.pos+1 < len(d.data) && d.data[d.pos+1] == ')' {
			d.pos += 2
			return nil, nil
		}
		// tuple
		d.pos++
		var arr []any
		first := true
		for {
			d.skipWhitespaceAndComments()
			if d.pos >= len(d.data) || d.data[d.pos] == ')' {
				d.pos++
				break
			}
			if !first {
				if d.data[d.pos] == ',' {
					d.pos++
					d.skipWhitespaceAndComments()
					if d.pos < len(d.data) && d.data[d.pos] == ')' {
						d.pos++
						break
					}
				}
			}
			first = false
			v, err := d.parseAnyValue()
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	case b == 't' || b == 'f':
		v, err := d.parseBool()
		if err == nil {
			return v, nil
		}
		// fallthrough to string
		return d.parsePlainValue()
	case b == '-' && d.pos+1 < len(d.data) && d.data[d.pos+1] >= '0' && d.data[d.pos+1] <= '9':
		return d.parseNumberAny()
	case b >= '0' && b <= '9':
		return d.parseNumberAny()
	default:
		return d.parsePlainValue()
	}
}

func (d *decoder) parseNumberAny() (any, error) {
	start := d.pos
	if d.pos < len(d.data) && d.data[d.pos] == '-' {
		d.pos++
	}
	digitStart := d.pos
	for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
		d.pos++
	}
	if d.pos == digitStart {
		// '-' with no digits — fall back to plain string (e.g. "-foo").
		d.pos = start
		return d.parsePlainValue()
	}
	isFloat := false
	// ABNF: integer part already required ≥1 digit. The fractional part, if
	// present, also requires ≥1 digit; tokens like "5." fall through to
	// plain-string per the type-priority cascade.
	if d.pos < len(d.data) && d.data[d.pos] == '.' {
		dotPos := d.pos
		d.pos++
		fracStart := d.pos
		for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
			d.pos++
		}
		if d.pos == fracStart {
			// "5." → not a number; rewind to before '.' so atValueEnd can
			// reject it and we fall through to plain-string.
			d.pos = dotPos
		} else {
			isFloat = true
		}
	}
	// Scientific notation: ["e"/"E"] ["+"/"-"] 1*DIGIT — required exponent
	// digits per the extended ABNF.
	if d.pos < len(d.data) && (d.data[d.pos] == 'e' || d.data[d.pos] == 'E') {
		expMark := d.pos
		d.pos++
		if d.pos < len(d.data) && (d.data[d.pos] == '+' || d.data[d.pos] == '-') {
			d.pos++
		}
		expStart := d.pos
		for d.pos < len(d.data) && d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
			d.pos++
		}
		if d.pos == expStart {
			// "1e" / "1e+" → not a number; rewind so it falls through.
			d.pos = expMark
		} else {
			isFloat = true
		}
	}
	// SPEC §8.1: digits followed by non-delimiter chars (e.g. "123abc") are
	// a plain string, not a number error.
	if !d.atValueEnd() {
		d.pos = start
		return d.parsePlainValue()
	}
	s := unsafeString(d.data[start:d.pos])
	if isFloat {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, d.errorf("invalid number")
		}
		return v, nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, d.errorf("invalid number")
	}
	return v, nil
}

// ---------------------------------------------------------------------------
// Compound type parsing
// ---------------------------------------------------------------------------

func (d *decoder) unmarshalSlice(fv reflect.Value) error {
	if err := d.enter(); err != nil {
		return err
	}
	defer d.leave()
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) || d.data[d.pos] != '[' {
		return d.errorf("expected '['")
	}
	elemCount := countTopLevelEntries(d.data, d.pos, '[', ']')
	d.pos++

	elemType := fv.Type().Elem()
	if elemCount > 0 {
		var slice reflect.Value
		if fv.Cap() < elemCount {
			slice = reflect.MakeSlice(fv.Type(), elemCount, elemCount)
		} else {
			slice = fv.Slice(0, elemCount)
		}
		first := true
		decoded := 0
		for decoded < elemCount {
			d.skipWhitespaceAndComments()
			if d.pos >= len(d.data) || d.data[d.pos] == ']' {
				if d.pos < len(d.data) && d.data[d.pos] == ']' {
					d.pos++
				}
				break
			}
			if !first {
				if d.data[d.pos] == ',' {
					d.pos++
					d.skipWhitespaceAndComments()
					if d.pos < len(d.data) && d.data[d.pos] == ']' {
						d.pos++
						break
					}
				} else {
					break
				}
			}
			first = false

			elem := slice.Index(decoded)
			if elemType.Kind() == reflect.Struct {
				if err := d.unmarshalNestedStruct(elem); err != nil {
					return err
				}
			} else {
				if err := d.unmarshalValue(elem); err != nil {
					return err
				}
			}
			decoded++
		}
		d.skipWhitespaceAndComments()
		if d.pos < len(d.data) && d.data[d.pos] == ']' {
			d.pos++
		}
		fv.Set(slice.Slice(0, decoded))
		return nil
	}

	var slice reflect.Value
	if fv.Cap() > 0 {
		slice = fv.Slice(0, 0)
	} else {
		slice = reflect.MakeSlice(fv.Type(), 0, 4)
	}

	first := true
	for {
		d.skipWhitespaceAndComments()
		if d.pos >= len(d.data) || d.data[d.pos] == ']' {
			d.pos++
			break
		}
		if !first {
			if d.data[d.pos] == ',' {
				d.pos++
				d.skipWhitespaceAndComments()
				if d.pos < len(d.data) && d.data[d.pos] == ']' {
					d.pos++
					break
				}
			} else {
				break
			}
		}
		first = false

		elem := reflect.New(elemType).Elem()
		// If element is a struct inside an array, it needs nested parsing
		if elemType.Kind() == reflect.Struct {
			if err := d.unmarshalNestedStruct(elem); err != nil {
				return err
			}
		} else {
			if err := d.unmarshalValue(elem); err != nil {
				return err
			}
		}
		slice = reflect.Append(slice, elem)
	}

	fv.Set(slice)
	return nil
}

func (d *decoder) unmarshalNestedStruct(fv reflect.Value) error {
	if err := d.enter(); err != nil {
		return err
	}
	defer d.leave()
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) {
		return d.errorf("unexpected EOF")
	}

	// Check for inline schema: {field1,field2,...}:(val1,val2,...)
	if d.data[d.pos] == '{' {
		si := getStructInfo(fv.Type())
		fields, schemaKey, err := d.parseSchema()
		if err != nil {
			return err
		}
		d.skipWhitespaceAndComments()
		if d.pos >= len(d.data) || d.data[d.pos] != ':' {
			return d.errorf("expected ':'")
		}
		d.pos++
		d.skipWhitespaceAndComments()
		fieldMap := buildFieldMapCached(si, fields, schemaKey)
		return d.unmarshalTuple(fv, si, fieldMap)
	}

	// Positional tuple: (val1,val2,...)
	if d.data[d.pos] == '(' {
		si := getStructInfo(fv.Type())
		return d.unmarshalTuple(fv, si, si.identityFieldMap)
	}

	return d.errorf("expected '{' or '(' for struct")
}

// ---------------------------------------------------------------------------
// skipValue — skip any value without allocating
// ---------------------------------------------------------------------------

func (d *decoder) skipValue() error {
	d.skipWhitespaceAndComments()
	if d.pos >= len(d.data) {
		return nil
	}

	switch d.data[d.pos] {
	case '"':
		d.pos++
		for d.pos < len(d.data) {
			if d.data[d.pos] == '"' {
				d.pos++
				return nil
			}
			if d.data[d.pos] == '\\' {
				d.pos++
			}
			d.pos++
		}
		return d.errorf("unclosed string")
	case '(':
		return d.skipBalanced('(', ')')
	case '[':
		return d.skipBalanced('[', ']')
	default:
		for d.pos < len(d.data) {
			b := d.data[d.pos]
			if b == ',' || b == ')' || b == ']' {
				return nil
			}
			d.pos++
		}
		return nil
	}
}

// ---------------------------------------------------------------------------
// unsafe helpers
// ---------------------------------------------------------------------------

func unsafeBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
