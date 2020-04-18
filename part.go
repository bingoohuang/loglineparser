package loglineparser

import (
	"github.com/rivo/uniseg"
)

// PartSplitter 表示日志行分割器
type PartSplitter interface {
	// ParseParts 解析成各个部分
	Parse(s string) []string
}

// bracketPartSplitter 定义了[]分割器
type bracketPartSplitter struct {
	emptyPlaceholder string
}

// NewBracketPartSplitter creates a new instance of PartSplitter.
func NewBracketPartSplitter(emptyPlaceholder string) PartSplitter {
	return &bracketPartSplitter{emptyPlaceholder: emptyPlaceholder}
}

// LoadLine 初始化
func (b bracketPartSplitter) Parse(s string) []string {
	gr := uniseg.NewGraphemes(s)
	reserved := ""
	parts := make([]string, 0)

	var (
		p  string
		ok bool
	)

	for {
		reserved, p, ok = b.next(gr, reserved)
		if !ok {
			break
		}

		parts = append(parts, p)
	}

	return parts
}

// nolint gocognit
// next 返回（reserved, part, ok)
func (b bracketPartSplitter) next(gr *uniseg.Graphemes, reserved string) (string, string, bool) {
	last := ""
	word := ""
	found := false
	maybeEnd := false

	reserveUsed := false
	s := reserved

	for {
		if !reserveUsed {
			reserveUsed = true
			if s != "" {
				goto PROCESS
			}
		}

		if !gr.Next() {
			break
		}
		s = gr.Str()

	PROCESS:
		if maybeEnd {
			if IsBlank(s) || !IsAlphanumeric(s) {
				if word == b.emptyPlaceholder {
					word = ""
				}
				return s, word, true
			}

			maybeEnd = false
			word += "]" + s
			goto LAST
		}

		if found {
			if s == "]" {
				maybeEnd = true
			} else {
				word += s
			}

			goto LAST
		}

		if s == "[" && IsBlank(last) {
			found = true
		}

	LAST:
		last = s
	}

	if maybeEnd {
		return "", word, true
	}

	return "", "", false
}
