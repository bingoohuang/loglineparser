package loglineparser

// Unmarshaler defines the interface for unmarshaling from a string.
type Unmarshaler interface {
	Unmarshal(string) error
}
