package loglineparser

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

func ParseLogLine(line string, result interface{}) error {
	parts := MakeBracketPartSplitter().Parse(line)
	subSplitter := MakeSubSplitter(",", "-")
	return ParseLog(parts, result, subSplitter)
}

var fieldsCache StructFieldsCache

func ParseLog(parts []string, result interface{}, subSplitter PartSplitter) error {
	v, err := CheckStructPtr(result)
	if err != nil {
		return err
	}

	structFields := fieldsCache.CachedStructFields(v.Type(), func(fieldIndex int, f reflect.StructField) interface{} {
		tag, _ := f.Tag.Lookup("llp")
		if !f.Anonymous && (tag == "" || tag == "-") {
			return nil
		}

		partIndex, subIndex := parseTwoInts(tag, -1)
		return structField{
			FieldIndex: fieldIndex,
			PartIndex:  partIndex,
			SubIndex:   subIndex,
			Kind:       f.Type.Kind(),
			Type:       f.Type,
			PtrType:    reflect.PtrTo(f.Type),
			Anonymous:  f.Anonymous,
		}
	}).([]structField)

	for _, sf := range structFields {
		err := fillField(parts, subSplitter, sf, v.Field(sf.FieldIndex))
		if err != nil {
			return err
		}
	}

	return nil
}

// structField 表示一个struct的字段属性
type structField struct {
	FieldIndex int
	PartIndex  int
	SubIndex   int
	Kind       reflect.Kind
	Type       reflect.Type
	PtrType    reflect.Type
	Anonymous  bool
}

func parseTwoInts(tag string, defaultValue int) (int, int) {
	s0, s1 := Split2(tag, ".")
	return ParseInt(s0, defaultValue), ParseInt(s1, defaultValue)
}

var unmarsherType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func fillField(parts []string, subSplitter PartSplitter, sf structField, f reflect.Value) error {
	if sf.Kind == reflect.Struct && sf.Anonymous {
		fv := reflect.New(f.Type()).Interface()
		err := ParseLog(parts, fv, subSplitter)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(fv).Elem())
		return nil
	}

	part := parsePart(sf, parts)
	if part == "" {
		return nil
	}

	sub := parseSub(subSplitter, part, sf)
	if sub == "" {
		return nil
	}

	if AssignBasicValue(f, sub) {
		return nil
	}

	var fv interface{}
	var err error

	switch f.Interface().(type) {
	case time.Time:
		v := ParseTime(sub)
		fv = &v
	default:
		if sf.Kind == reflect.Map {
			fv = reflect.New(f.Type()).Interface()
			err = json.Unmarshal([]byte(sub), fv)
		} else if sf.PtrType.Implements(unmarsherType) {
			fv = reflect.New(f.Type()).Interface()
			err = fv.(Unmarshaler).Unmarshal(sub)
		}
	}

	if err != nil {
		return err
	}
	if fv != nil {
		f.Set(reflect.ValueOf(fv).Elem())
		return nil
	}

	return errors.New(sf.Kind.String() + " is not supported")
}

func parseSub(subSplitter PartSplitter, part string, sf structField) string {
	if sf.SubIndex < 0 {
		return part
	}

	subs := subSplitter.Parse(part)
	if sf.SubIndex < len(subs) {
		return subs[sf.SubIndex]
	}

	return ""
}

func parsePart(sf structField, parts []string) string {
	if sf.PartIndex < 0 {
		return ""
	}

	part := ""
	if sf.PartIndex < len(parts) {
		part = parts[sf.PartIndex]
	}
	if part == "-" {
		part = ""
	}
	return part
}
