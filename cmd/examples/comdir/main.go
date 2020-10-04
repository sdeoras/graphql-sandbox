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
	log.Init(zapcore.ErrorLevel)
	comdirutil.Init()

	h, err := auth.NewHandler(&auth.Config{
		PublicKey: "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFcWNoUlMvMGd1RURGMmxBK0ZxdC8rWG9IYXVUcgorSHZBUWtXMW1iVndnRVNHNmdFUXZtaDNiVjN2cThNeDRxaG5IdjVwRm5IeEp0eUtnQVdEb3pFMmxnPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==",
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
	})
	if err != nil {
		log.Logger().Fatal(err.Error())
		return
	}

	http.Handle("/", h)

	fmt.Println("Now server is running on port 8081")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8081/")
	fmt.Println("JWT Token 1:")
	fmt.Println("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsInBhZCI6Ii4uLi4uIn0K.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZW1haWwiOiJqb2huLmRvZUBnb29nbGUuY29tIn0K.-SK5uwI3qeDQilqAEqwzRpb6aE5uWpgcWTPoXb4LC6sZuB20e0NfSmYKjQNMnRrfWckKOg-gnNUMI0FSrUB5sw")
	fmt.Println()
	fmt.Println("JWT Token 2:")
	fmt.Println("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsInBhZCI6Ii4uLi4uIn0K.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkFsaWNlIERlZSIsImVtYWlsIjoiYWxpY2UuZGVlQGFwcGxlLmNvbSJ9Cg==.esAGhPqy9bjsH5IBAY9OBR_BpxWhlgLxHbB-1e9EV3U8KyRrFt5jFngkAqTPQcBwVDd-OJHFeBIeKyatFXxmaw")
	fmt.Println()
	fmt.Println("Query:")
	fmt.Println(q)
	_ = http.ListenAndServe(":8081", nil)
}

var q = `
query jonDoeLogin {
  login(username:"jon.doe@google.com", password:"abcd") {
    jwt
  }
}

query aliceDeeLogin {
  login(username:"alice.dee@apple.com", password:"abcd") {
    jwt
  }
}

query employee1 {
  employee(id: "id-1") {
    name
  }
}

query employee2 {
  employee(id: "id-2") {
    name
  }
}

mutation m {
  employee(id: "id-1", name: "alice") {
    name
  }
}
`
