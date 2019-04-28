package loglineparser

type logPartSplitter struct {
	sep              string
	emptyPlaceholder string
}

func MakeSubSplitter(sep, emptyPlaceholder string) PartSplitter {
	return &logPartSplitter{sep: sep, emptyPlaceholder: emptyPlaceholder}
}

func (l *logPartSplitter) Parse(part string) []string {
	subs := SplitN(part, l.sep, true, false)
	for i, p := range subs {
		if p == "-" {
			subs[i] = ""
		}
	}

	return subs
}
