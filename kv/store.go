package kv

// Store holds values mapped to a key.
type Store interface {

	// Set sets the value for the key
	Set(key, value string)

	// Get retrieves the value associated with key. The ok result indicates whether value was found in the store. If none found the value will be a zero length string
	Get(key string) (value string, ok bool)
}
