package api

import "net/http"

type HttpRequestReader interface {
	ReadFrom(r *http.Request) error
}
