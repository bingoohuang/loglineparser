package loglineparser

import (
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/spf13/cast"
)

// nolint gochecknoglobals
var reg = regexp.MustCompile(`\d+\.\d+`)

// ParseTime parses s as a time.Time.
func ParseTime(s interface{}) time.Time {
	if vt, ok := s.(string); ok {
		if reg.MatchString(vt) {
			sec, millis := parseTwoInts(vt)
			return time.Unix(int64(*sec), int64(*millis*1000000)) // nolint gomnd
		}
	}

	t, _ := dateparse.ParseAny(cast.ToString(s))

	return t
}
