package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/sdeoras/graphql/pkg/testutil/comdirutil"
)

func main() {
	http.Handle("/",
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
	)

	fmt.Println("Now server is running on port 8081")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}