package asun

import (
	"encoding/binary"
	"math"
	"reflect"
	"sync"
)

// appendUvarint appends v as an LEB128 unsigned varint.
func appendUvarint(buf []byte, v uint64) []byte {
	for v >= 0x80 {
		buf = append(buf, byte(v)|0x80)
		v >>= 7
	}
	return append(buf, byte(v))
}

// zigzagEncode maps a signed integer to an unsigned one so that small
// magnitudes (positive or negative) encode into few varint bytes.
func zigzagEncode(v int64) uint64 {
	return uint64((v << 1) ^ (v >> 63))
}

// appendIvarint appends v as zigzag + LEB128 signed varint.
func appendIvarint(buf []byte, v int64) []byte {
	return appendUvarint(buf, zigzagEncode(v))
}

var binBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 128)
		return &b
	},
}

func getBinBuf() *[]byte {
	bp := binBufPool.Get().(*[]byte)
	*bp = (*bp)[:0]
	return bp
}

func putBinBuf(bp *[]byte) {
	if cap(*bp) <= 1<<16 {
		binBufPool.Put(bp)
	}
}

// EncodeBinary serializes a Go value to ASUN-BIN format.
func EncodeBinary(v any) ([]byte, error) {
	if v == nil {
		return nil, &MarshalError{Message: "cannot marshal nil"}
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, &MarshalError{Message: "cannot marshal nil pointer"}
		}
		rv = rv.Elem()
	}
	if err := ensureNoMapType(rv.Type()); err != nil {
		return nil, err
	}

	bp := getBinBuf()
	buf := *bp
	var err error
	buf, err = marshalBinValue(buf, rv)
	if err != nil {
		*bp = buf
		putBinBuf(bp)
		return nil, err
	}
	result := make([]byte, len(buf))
	copy(result, buf)
	*bp = buf
	putBinBuf(bp)
	return result, nil
}

func marshalBinValue(buf []byte, rv reflect.Value) ([]byte, error) {
	switch rv.Kind() {
	case reflect.Bool:
		if rv.Bool() {
			return append(buf, 1), nil
		}
		return append(buf, 0), nil
	case reflect.Int8:
		return append(buf, byte(rv.Int())), nil
	case reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		buf = appendIvarint(buf, rv.Int())
		return buf, nil
	case reflect.Uint8:
		return append(buf, byte(rv.Uint())), nil
	case reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		buf = appendUvarint(buf, rv.Uint())
		return buf, nil
	case reflect.Float32:
		buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(float32(rv.Float())))
		return buf, nil
	case reflect.Float64:
		buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(rv.Float()))
		return buf, nil
	case reflect.String:
		s := rv.String()
		buf = appendUvarint(buf, uint64(len(s)))
		buf = append(buf, s...)
		return buf, nil
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			b := rv.Bytes()
			buf = appendUvarint(buf, uint64(len(b)))
			buf = append(buf, b...)
			return buf, nil
		}
		n := rv.Len()
		buf = appendUvarint(buf, uint64(n))
		for i := 0; i < n; i++ {
			var err error
			buf, err = marshalBinValue(buf, rv.Index(i))
			if err != nil {
				return buf, err
			}
		}
		return buf, nil
	case reflect.Array:
		n := rv.Len()
		buf = appendUvarint(buf, uint64(n))
		for i := 0; i < n; i++ {
			var err error
			buf, err = marshalBinValue(buf, rv.Index(i))
			if err != nil {
				return buf, err
			}
		}
		return buf, nil
	case reflect.Map:
		return buf, errMapFieldsUnsupported
	case reflect.Struct:
		si := getStructInfo(rv.Type())
		for _, f := range si.fields {
			fv := rv.FieldByIndex(f.index)
			var err error
			buf, err = marshalBinValue(buf, fv)
			if err != nil {
				return buf, err
			}
		}
		return buf, nil
	case reflect.Ptr:
		if rv.IsNil() {
			return append(buf, 0), nil
		}
		buf = append(buf, 1)
		return marshalBinValue(buf, rv.Elem())
	case reflect.Interface:
		if rv.IsNil() {
			return append(buf, 0), nil
		}
		buf = append(buf, 1)
		return marshalBinValue(buf, rv.Elem())
	default:
		return buf, nil
	}
}
