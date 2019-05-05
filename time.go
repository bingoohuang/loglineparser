package loglineparser

import (
	"github.com/araddon/dateparse"
	"github.com/spf13/cast"
	"regexp"
	"time"
)

var reg = regexp.MustCompile(`\d+\.\d+`)

func ParseTime(s interface{}) time.Time {
	switch vt := s.(type) {
	case string:
		if reg.MatchString(vt) {
			sec, millis := parseTwoInts(vt, 0)
			return time.Unix(int64(sec), int64(millis*1000000))
		}
	}

	t, _ := dateparse.ParseAny(cast.ToString(s))
	return t
}
