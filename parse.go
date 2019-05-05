package loglineparser

import (
	"encoding/json"
	"errors"
	"github.com/modern-go/reflect2"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type LogLineParser struct {
	FieldsCache  StructFieldsCache
	StructType   reflect2.Type
	PartSplitter PartSplitter
	SubSplitter  PartSplitter
}

func NewLogLineParser(structPointerOrName interface{}) *LogLineParser {
	structTypePtr := reflect.TypeOf(structPointerOrName)
	kind := structTypePtr.Kind()
	if kind == reflect.String {
		return &LogLineParser{
			StructType:   reflect2.TypeByName(structPointerOrName.(string)),
			PartSplitter: NewBracketPartSplitter("-"),
			SubSplitter:  NewSubSplitter(",", "-"),
		}
	}

	if kind != reflect.Ptr {
		logrus.Panicf("non struct ptr %v", structTypePtr)
	}

	elem := structTypePtr.Elem()
	if elem.Kind() != reflect.Struct {
		logrus.Panicf("non struct ptr %v", structTypePtr)
	}

	return NewLogLineParser(elem.String())
}

func (l *LogLineParser) Parse(line string) (interface{}, error) {
	p := l.StructType.New()
	err := l.parse(l.PartSplitter.Parse(line), p)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(p).Elem().Interface(), nil
}

func createStructField(fieldIndex int, f reflect.StructField) interface{} {
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
}

func (l *LogLineParser) parse(parts []string, result interface{}) error {
	v := reflect.ValueOf(result).Elem()
	structFields := l.FieldsCache.CachedStructFields(v.Type(), createStructField).([]structField)

	for _, sf := range structFields {
		err := l.fillField(parts, sf, v.Field(sf.FieldIndex))
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

func (l *LogLineParser) fillField(parts []string, sf structField, f reflect.Value) error {
	if sf.Kind == reflect.Struct && sf.Anonymous {
		fv := reflect.New(f.Type()).Interface()
		err := l.parse(parts, fv)
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

	sub := l.parseSub(part, sf)
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

func (l *LogLineParser) parseSub(part string, sf structField) string {
	if sf.SubIndex < 0 {
		return part
	}

	subs := l.SubSplitter.Parse(part)
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

	return part
}
