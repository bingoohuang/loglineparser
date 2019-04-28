package loglineparser

type SubSplitter interface {
	Load(part, sep string)
	Len() int
	Subs() []string
	Sub(i int) string
}

type logPartSplitter struct {
	sep  string
	subs []string
}

func MakeSubSplitter() SubSplitter {
	return &logPartSplitter{}
}

func (l *logPartSplitter) Load(part, sep string) {
	subs := SplitN(part, sep, true, false)
	for i, p := range subs {
		if p == "-" {
			subs[i] = ""
		}
	}

	l.sep = sep
	l.subs = subs
}

func (l logPartSplitter) Subs() []string {
	return l.subs
}

func (l logPartSplitter) Len() int {
	return len(l.subs)
}

func (l logPartSplitter) Sub(i int) string {
	if i < len(l.subs) {
		return l.subs[i]
	}

	return ""
}
