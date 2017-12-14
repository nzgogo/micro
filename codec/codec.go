package codec

// Codec interface
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
