package main

import (
	"fmt"
	"github.com/graphql-go/handler"
	"github.com/sdeoras/graphql/pkg/log"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"github.com/sdeoras/graphql/pkg/testutil/comdirutil"
	"net/http"
)

func main() {
	http.Handle("/",
		auth.NewHandler(&auth.Config{
			Handler: handler.New(
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
			SkipCheck: true,
			Logger:    log.Logger(),
		}),
	)

	fmt.Println("Now server is running on port 8081")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8081/")
	_ = http.ListenAndServe(":8081", nil)
}
