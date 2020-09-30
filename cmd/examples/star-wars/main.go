package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/sdeoras/graphql/pkg/testutil/starwarutil"
)

func main() {
	http.Handle("/",
		handler.New(
			&handler.Config{
				Schema:           &starwarutil.Schema,
				Pretty:           true,
				GraphiQL:         true,
				Playground:       true,
				RootObjectFn:     nil,
				ResultCallbackFn: nil,
				FormatErrorFn:    nil,
			},
		),
	)

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/?query={hero{name}}'")
	http.ListenAndServe(":8080", nil)
}
