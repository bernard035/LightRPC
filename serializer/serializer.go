package serializer

// Serializer is interface, each serializer has Marshal and Unmarshal functions
type Serializer interface {
	Marshal(message interface{}) ([]byte, error)
	Unmarshal(data []byte, message interface{}) error
}
