package loglineparser

type Unmarshaler interface {
	Unmarshal(string) error
}
