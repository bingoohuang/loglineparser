package loglineparser

type logPartSplitter struct {
	sep              string
	emptyPlaceholder string
}

// NewSubSplitter creates a new PartSplitter.
func NewSubSplitter(sep, emptyPlaceholder string) PartSplitter {
	return &logPartSplitter{sep: sep, emptyPlaceholder: emptyPlaceholder}
}

// Parse parses the log parts.
func (l *logPartSplitter) Parse(s string) []string {
	subs := SplitN(s, l.sep, true, false)
	for i, p := range subs {
		if p == l.emptyPlaceholder {
			subs[i] = ""
		}
	}

	return subs
}
