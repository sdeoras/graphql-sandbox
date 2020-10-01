package main

import (
	"fmt"
	"github.com/graphql-go/handler"
	handler2 "github.com/sdeoras/graphql/pkg/rest/mw/auth/handler"
	"github.com/sdeoras/graphql/pkg/testutil/comdirutil"
	"net/http"
)

func main() {
	http.Handle("/",
		handler2.New(
			handler.New(
				&handler.Config{
					Schema:           &comdirutil.Schema,
					Pretty:           true,
					GraphiQL:         true,
					Playground:       true,
					RootObjectFn:     nil,
					ResultCallbackFn: nil,
					FormatErrorFn:    nil,
				},
			),
			true,
		),
	)

	fmt.Println("Now server is running on port 8081")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
