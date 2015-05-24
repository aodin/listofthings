package config

// Metadata holds arbitrary strings as key - value pairs
type Metadata map[string]string

// Get returns the value of the given key. If the key does not exist in the
// metadata, a blank string will be returned
func (m Metadata) Get(key string) string {
	return m[key]
}

// Has returns true if the metadata contains the key. Keys with blank values
// will return true.
func (m Metadata) Has(key string) (exists bool) {
	_, exists = m[key]
	return
}

// Keys returns all keys of the metadata
func (m Metadata) Keys() []string {
	keys := make([]string, len(m))
	var i int
	for key, _ := range m {
		keys[i] = key
		i += 1
	}
	return keys
}

// Values returns all values of the metadata
func (m Metadata) Values() []string {
	values := make([]string, len(m))
	var i int
	for _, value := range m {
		values[i] = value
		i += 1
	}
	return values
}
