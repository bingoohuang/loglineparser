package loglineparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/modern-go/reflect2"
)

// LogLineParser defines the struct representing parser for log line.
type LogLineParser struct {
	FieldsCache  StructFieldsCache
	StructType   reflect2.Type
	Ptr          bool
	PartSplitter PartSplitter
	SubSplitter  PartSplitter
}

// New creates a new LogLineParser.
func New(structPointerOrName interface{}) (*LogLineParser, error) {
	structTypePtr := reflect.TypeOf(structPointerOrName)

	switch kind := structTypePtr.Kind(); kind {
	case reflect.String:
		name := structPointerOrName.(string)
		ptr := strings.HasPrefix(name, "*")
		if ptr {
			name = name[1:]
		}
		return &LogLineParser{
			StructType:   reflect2.TypeByName(name),
			Ptr:          ptr,
			PartSplitter: NewBracketPartSplitter("-"),
			SubSplitter:  NewSubSplitter(",", "-"),
		}, nil
	case reflect.Struct:
		return &LogLineParser{
			StructType:   reflect2.Type2(structTypePtr),
			Ptr:          false,
			PartSplitter: NewBracketPartSplitter("-"),
			SubSplitter:  NewSubSplitter(",", "-"),
		}, nil
	case reflect.Ptr:
		elem := structTypePtr.Elem()
		if elem.Kind() != reflect.Struct {
			return nil, fmt.Errorf("non struct ptr %v", structTypePtr)
		}

		return &LogLineParser{
			StructType:   reflect2.Type2(elem),
			Ptr:          true,
			PartSplitter: NewBracketPartSplitter("-"),
			SubSplitter:  NewSubSplitter(",", "-"),
		}, nil
	}

	return nil, fmt.Errorf("non struct ptr %v", structTypePtr)
}

// Parse parses a line string.
func (l *LogLineParser) Parse(line string) (interface{}, error) {
	p := l.StructType.New()

	parts := l.PartSplitter.Parse(line)
	err := l.parseParts(line, parts, p)
	if err != nil {
		return nil, err
	}

	if l.Ptr {
		return p, nil
	}

	return reflect.ValueOf(p).Elem().Interface(), nil
}

func createStructField(fieldIndex int, f reflect.StructField) interface{} {
	tag, _ := f.Tag.Lookup("llp")
	if !f.Anonymous && (tag == "" || tag == "-") {
		return nil
	}

	sf := structField{
		FieldIndex: fieldIndex,
		Kind:       f.Type.Kind(),
		Type:       f.Type,
		PtrType:    reflect.PtrTo(f.Type),
		Anonymous:  f.Anonymous,
	}

	if tag == "reg" {
		group := 1
		var err error
		if groupTag := f.Tag.Get("group"); groupTag != "" {
			if group, err = strconv.Atoi(groupTag); err != nil {
				log.Fatalf("group %s is not a valid number", groupTag)
			}
		}

		regTag := f.Tag.Get("reg")
		sf.Regexp = regexp.MustCompile(regTag)
		sf.RegexpGroup = group
	} else {
		sf.PartIndex, sf.SubIndex = parseTwoInts(tag, -1)
	}

	return sf
}

func (l *LogLineParser) parseParts(line string, parts []string, result interface{}) error {
	v := reflect.ValueOf(result).Elem()
	structFields := l.FieldsCache.CachedStructFields(v.Type(), createStructField).([]structField)

	for _, sf := range structFields {
		err := l.fillField(line, parts, sf, v.Field(sf.FieldIndex))
		if err != nil {
			return err
		}
	}

	return nil
}

// structField 表示一个struct的字段属性
type structField struct {
	FieldIndex  int
	PartIndex   int
	SubIndex    int
	Kind        reflect.Kind
	Type        reflect.Type
	PtrType     reflect.Type
	Anonymous   bool
	Regexp      *regexp.Regexp
	RegexpGroup int
}

func parseTwoInts(tag string, defaultValue int) (int, int) {
	s0, s1 := Split2(tag, ".")
	return ParseInt(s0, defaultValue), ParseInt(s1, defaultValue)
}

// nolint gochecknoglobals
var unmarsherType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func (l *LogLineParser) fillField(line string, parts []string, sf structField, f reflect.Value) error {
	if sf.Kind == reflect.Struct && sf.Anonymous {
		fv := reflect.New(f.Type()).Interface()
		if err := l.parseParts(line, parts, fv); err != nil {
			return err
		}

		f.Set(reflect.ValueOf(fv).Elem())

		return nil
	}

	var sub string

	if sf.Regexp != nil {
		subs := sf.Regexp.FindStringSubmatch(line)
		if sf.RegexpGroup < len(subs) {
			sub = subs[sf.RegexpGroup]
		}
	} else {
		part := parsePart(sf, parts)
		if part == "" {
			return nil
		}

		sub = l.parseSub(part, sf)
	}

	if sub == "" {
		return nil
	}

	if AssignBasicValue(f, sub) {
		return nil
	}

	var (
		fv  interface{}
		err error
	)

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
