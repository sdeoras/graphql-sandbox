package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/sdeoras/graphql/pkg/log"
	"github.com/sdeoras/graphql/pkg/rest/mw/auth"
	"github.com/sdeoras/graphql/pkg/testutil/comdirutil"
	"go.uber.org/zap/zapcore"
)

func main() {
	log.Init(zapcore.DebugLevel)
	comdirutil.Init()

	http.Handle("/",
		auth.NewHandler(&auth.Config{
			Handler: handler.New(
				&handler.Config{
					Schema:           &comdirutil.Schema,
					Pretty:           true,
					GraphiQL:         false,
					Playground:       true,
					RootObjectFn:     nil,
					ResultCallbackFn: nil,
					FormatErrorFn:    nil,
				},
			),
			Logger: log.Logger(),
		}),
	)

	fmt.Println("Now server is running on port 8081")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8081/")
	_ = http.ListenAndServe(":8081", nil)
}
