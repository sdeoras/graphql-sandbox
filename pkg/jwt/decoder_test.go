package jwt

import (
	"fmt"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	token := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsInBhZCI6Ii4uLi4uIn0K.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOiJqb2huLmRvZUBnb29nbGUuY29tIn0K.-SK5uwI3qeDQilqAEqwzRpb6aE5uWpgcWTPoXb4LC6sZuB20e0NfSmYKjQNMnRrfWckKOg-gnNUMI0FSrUB5sw"

	decoder, err := NewDecoder(publicKey)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := decoder.Decode(token)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range claims {
		fmt.Println(key, value)
	}
}
