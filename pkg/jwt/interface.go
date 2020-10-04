package jwt

type Decoder interface {
	Decode(token string) (map[string]interface{}, error)
}

func NewDecoder(keyB64 string) (Decoder, error) {
	return newDecoder(keyB64)
}
