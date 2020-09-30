package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GraphQlRequest struct {
	Query string `json:"query"`
}

func (g *GraphQlRequest) ReadFrom(r *http.Request) error {
	if r == nil {
		return fmt.Errorf("http request is nil and therefore invalid")
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading http request: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, g); err != nil {
		return fmt.Errorf("error unmarshaling http reqeust: %w", err)
	}

	return nil
}
