package cache

// Entry represents a key/value pair.
type Entry struct {
	Key   []byte
	Value []byte
}

type Iterator interface {
	Next() *Entry
}
