package loglineparser

import (
	"github.com/rivo/uniseg"
)

// PartSplitter 表示日志行分割器
type PartSplitter interface {
	// LoadLine 初始化
	LoadLine(line string)
	// ParseParts 解析成各个部分
	ParseParts() []string
}

// bracketPartSplitter 定义了[]分割器
type bracketPartSplitter struct {
	gr       *uniseg.Graphemes
	index    int
	offset   int
	reserved string

	parts  []string
	parsed bool
}

func MakeBracketPartSplitter() PartSplitter {
	return &bracketPartSplitter{}
}

// LoadLine 初始化
func (b *bracketPartSplitter) LoadLine(line string) {
	b.gr = uniseg.NewGraphemes(line)
	b.index = 0
	b.offset = 0
	b.reserved = ""
	b.parts = make([]string, 0)
	b.parsed = false
}

func (b *bracketPartSplitter) ParseParts() []string {
	if b.parsed {
		return b.parts
	}
	for {
		p, pi := b.next()
		if pi < 0 {
			break
		}

		b.parts = append(b.parts, p)
	}

	b.parsed = true
	return b.parts
}

// Next 返回下一个分隔部分，返回分隔好的部分part和当前part的索引.
// partIndex == -1时，表示分隔已经完成.

func (b *bracketPartSplitter) next() (part string, partIndex int) {
	last := ""
	word := ""
	found := false
	maybeEnd := false

	reserveUsed := false
	s := b.reserved
	b.reserved = ""

	for {
		if !reserveUsed {
			reserveUsed = true
			if s != "" {
				goto PROCESS
			}
		}

		if !b.gr.Next() {
			break
		}
		s = b.gr.Str()

	PROCESS:
		if maybeEnd {
			if IsBlank(s) || !IsAlphanumeric(s) {
				partIndex = b.index
				b.index++
				b.reserved = s
				return word, partIndex
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

	return "", -1
}
