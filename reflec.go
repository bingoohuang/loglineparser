package loglineparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"sync"
	"time"
)

func ParseLogLine(line string, result interface{}) error {
	partSplitter := MakeBracketPartSplitter()
	partSplitter.LoadLine(line)

	subSplitter := MakeSubSplitter()

	return ParseLog(result, partSplitter, subSplitter)
}

func ParseLog(result interface{}, lineSplitter PartSplitter, partSplitter SubSplitter) error {
	structTypePtr := reflect.TypeOf(result)
	if structTypePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("non struct ptr %v", structTypePtr)
	}
	v := reflect.ValueOf(result).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("non struct ptr %v", structTypePtr)
	}

	structFields := cachedStructFields(v.Type())
	for i, sf := range structFields {
		err := fillField(lineSplitter, partSplitter, sf, v.Field(i))
		if err != nil {
			return err
		}
	}

	return nil
}

var fieldCache sync.Map // map[reflect.Type][]field

// StructField 表示一个struct的字段属性
type StructField struct {
	PartIndex int
	SubIndex  int
	Kind      reflect.Kind
	Type      reflect.Type
	PtrType   reflect.Type
	Anonymous bool
}

// CachedStructFields caches fields of struct type
func cachedStructFields(t reflect.Type) []StructField {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]StructField)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]StructField)
}

func typeFields(t reflect.Type) []StructField {
	ff := t.NumField()
	fields := make([]StructField, ff)

	for fi := 0; fi < ff; fi++ {
		f := t.Field(fi)

		tag, _ := f.Tag.Lookup("llp")
		if tag == "" || tag == "-" {
			continue
		}

		partIndex, subIndex := parseTwoInts(tag, -1)
		fields[fi] = StructField{
			PartIndex: partIndex,
			SubIndex:  subIndex,
			Kind:      f.Type.Kind(),
			Type:      f.Type,
			PtrType:   reflect.PtrTo(f.Type),
			Anonymous: f.Anonymous,
		}
	}

	return fields
}

func parseTwoInts(tag string, defaultValue int) (int, int) {
	s0, s1 := Split2(tag, ".")
	return parseInt(s0, defaultValue), parseInt(s1, defaultValue)
}

func parseInt(s string, defaultValue int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	return defaultValue
}

var unmarsherType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func fillField(lineSplitter PartSplitter, partSplitter SubSplitter, sf StructField, f reflect.Value) error {
	parts := lineSplitter.ParseParts()
	part := ""
	if sf.PartIndex >= 0 && sf.PartIndex < len(parts) {
		part = parts[sf.PartIndex]
	}
	if part == "-" {
		part = ""
	}

	partSplitter.Load(part, ",")
	subs := partSplitter.Subs()
	sub := ""
	if sf.SubIndex < 0 {
		sub = part
	} else if sf.SubIndex < len(subs) {
		sub = subs[sf.SubIndex]
	}

	if setFieldValue(sf, f, sub) {
		return nil
	}

	if sf.Kind == reflect.Struct && sf.Anonymous {
		fv := reflect.New(f.Type()).Interface()
		err := ParseLog(fv, lineSplitter, partSplitter)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(fv).Elem())
		return nil
	}

	if sf.Kind == reflect.Map {
		fv := reflect.New(f.Type()).Interface()
		err := json.Unmarshal([]byte(sub), fv)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(fv).Elem())
		return nil
	}

	if sf.PtrType.Implements(unmarsherType) {
		fv := reflect.New(f.Type()).Interface()
		err := fv.(Unmarshaler).Unmarshal(sub)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(fv).Elem())
		return nil
	}

	return errors.New(sf.Kind.String() + " is not supported")
}

func setFieldValue(sf StructField, f reflect.Value, v interface{}) bool {
	vt := reflect.TypeOf(v)
	if vt == f.Type() {
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
	case string:
		f.SetString(cast.ToString(v))
	case time.Time:
		t := ParseTime(v)
		f.Set(reflect.ValueOf(t))
	default:
		return false
	}

	return true
}
