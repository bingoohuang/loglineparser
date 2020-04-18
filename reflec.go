package loglineparser

import (
	"reflect"
	"sync"

	"github.com/spf13/cast"
)

// StructFieldsCache defines the cache of struct fields.
type StructFieldsCache struct {
	fieldsCache sync.Map // map[reflect.Type][]field
}

// CachedStructFields caches fields of struct type
func (s *StructFieldsCache) CachedStructFields(t reflect.Type,
	fn func(fieldIndex int, f reflect.StructField) interface{}) interface{} {
	if f, ok := s.fieldsCache.Load(t); ok {
		return f
	}

	f, _ := s.fieldsCache.LoadOrStore(t, typeFields(t, fn))

	return f
}

func typeFields(t reflect.Type, fn func(fieldIndex int, f reflect.StructField) interface{}) interface{} {
	var fields reflect.Value

	for fi := 0; fi < t.NumField(); fi++ {
		f := t.Field(fi)
		field := fn(fi, f)

		if field == nil {
			continue
		}

		fv := reflect.ValueOf(field)

		if !fields.IsValid() {
			fields = reflect.MakeSlice(reflect.SliceOf(fv.Type()), 0, t.NumField())
		}

		fields = reflect.Append(fields, fv)
	}

	return fields.Interface()
}

// AssignBasicValue assigns v as basic types lie string/(u)int(8/16/32/64) or its type self.
func AssignBasicValue(f reflect.Value, v interface{}) bool {
	if reflect.TypeOf(v) == f.Type() {
		f.Set(reflect.ValueOf(v))
		return true
	}

	switch f.Interface().(type) {
	case bool:
		f.SetBool(cast.ToBool(v))
	case float32, float64:
		f.SetFloat(cast.ToFloat64(v))
	case int8, int16, int, int32, int64:
		f.SetInt(cast.ToInt64(v))
	case uint8, uint16, uint, uint32, uint64:
		f.SetUint(cast.ToUint64(v))
	case string:
		f.SetString(cast.ToString(v))
	default:
		return false
	}

	return true
}
