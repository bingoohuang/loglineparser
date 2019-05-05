package loglineparser

type logPartSplitter struct {
	sep              string
	emptyPlaceholder string
}

func NewSubSplitter(sep, emptyPlaceholder string) PartSplitter {
	return &logPartSplitter{sep: sep, emptyPlaceholder: emptyPlaceholder}
}

func (l *logPartSplitter) Parse(s string) []string {
	subs := SplitN(s, l.sep, true, false)
	for i, p := range subs {
		if p == l.emptyPlaceholder {
			subs[i] = ""
		}
	}

	return subs
}
