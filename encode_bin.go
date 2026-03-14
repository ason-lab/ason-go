package ason

import (
	"encoding/binary"
	"math"
	"reflect"
	"sync"
)

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

// EncodeBinary serializes a Go value to ASON-BIN format.
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
	case reflect.Int16:
		buf = binary.LittleEndian.AppendUint16(buf, uint16(rv.Int()))
		return buf, nil
	case reflect.Int32:
		buf = binary.LittleEndian.AppendUint32(buf, uint32(rv.Int()))
		return buf, nil
	case reflect.Int, reflect.Int64:
		buf = binary.LittleEndian.AppendUint64(buf, uint64(rv.Int()))
		return buf, nil
	case reflect.Uint8:
		return append(buf, byte(rv.Uint())), nil
	case reflect.Uint16:
		buf = binary.LittleEndian.AppendUint16(buf, uint16(rv.Uint()))
		return buf, nil
	case reflect.Uint32:
		buf = binary.LittleEndian.AppendUint32(buf, uint32(rv.Uint()))
		return buf, nil
	case reflect.Uint, reflect.Uint64:
		buf = binary.LittleEndian.AppendUint64(buf, uint64(rv.Uint()))
		return buf, nil
	case reflect.Float32:
		buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(float32(rv.Float())))
		return buf, nil
	case reflect.Float64:
		buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(rv.Float()))
		return buf, nil
	case reflect.String:
		s := rv.String()
		buf = binary.LittleEndian.AppendUint32(buf, uint32(len(s)))
		buf = append(buf, s...)
		return buf, nil
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			b := rv.Bytes()
			buf = binary.LittleEndian.AppendUint32(buf, uint32(len(b)))
			buf = append(buf, b...)
			return buf, nil
		}
		n := rv.Len()
		buf = binary.LittleEndian.AppendUint32(buf, uint32(n))
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
		buf = binary.LittleEndian.AppendUint32(buf, uint32(n))
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
