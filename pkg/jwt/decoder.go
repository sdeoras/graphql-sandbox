package jwt

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var publicKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFcWNoUlMvMGd1RURGMmxBK0ZxdC8rWG9IYXVUcgorSHZBUWtXMW1iVndnRVNHNmdFUXZtaDNiVjN2cThNeDRxaG5IdjVwRm5IeEp0eUtnQVdEb3pFMmxnPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="

type decoder struct {
	pub *ecdsa.PublicKey
}

func newDecoder(keyB64 string) (*decoder, error) {
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, err
	}

	pub, err := jwt.ParseECPublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return &decoder{pub: pub}, nil
}

func (g *decoder) Decode(token string) (map[string]interface{}, error) {
	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return g.pub, nil
	})
	if err != nil {
		return nil, err
	}

	if t.Valid {
		return t.Claims.(jwt.MapClaims), nil
	}

	return nil, fmt.Errorf("token is not valid")
}
