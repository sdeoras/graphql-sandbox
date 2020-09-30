package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/sdeoras/graphql/pkg/log"
	"github.com/sdeoras/graphql/pkg/resolvers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log.Init(zapcore.DebugLevel)
	logger := log.Logger()
	defer logger.Sync()

	resolver, err := resolvers.NewResolver(
		&resolvers.Config{
			Logger: logger,
		},
	)
	if err != nil {
		logger.Fatal("fatal error", zap.String("err", err.Error()))
		return
	}

	const (
		hello     = "hello"
		rootQuery = "RootQuery"
	)

	f, err := resolver.Resolve(hello)
	if err != nil {
		logger.Fatal("fatal error", zap.String("err", err.Error()))
		return
	}

	// Schema
	fields := graphql.Fields{
		hello: &graphql.Field{
			Name:              "Hello",
			Type:              graphql.String,
			Args:              nil,
			Resolve:           f,
			DeprecationReason: "",
			Description:       "",
		},
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:        rootQuery,
			Interfaces:  nil,
			Fields:      fields,
			IsTypeOf:    nil,
			Description: "",
		}),
		Mutation:     nil,
		Subscription: nil,
		Types:        nil,
		Directives:   nil,
		Extensions:   nil,
	})
	if err != nil {
		msg := fmt.Errorf("failed to create graphql newschema:%w", err)
		logger.Fatal(msg.Error())
		return
	}

	// Query
	query := `
		{
			hello
		}
	`

	r := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		RootObject:     nil,
		VariableValues: nil,
		OperationName:  "",
		Context:        context.Background(),
	})

	if len(r.Errors) > 0 {
		msg := fmt.Sprintf("failed to execute graphql operation, errors: %+v", r.Errors)
		logger.Fatal(msg)
		return
	}

	jb, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		msg := fmt.Sprintf("failed to json marshal graphql response:%v", err)
		logger.Fatal(msg)
		return
	}

	fmt.Printf("%s\n", jb) // {"data":{"hello":"world"}}
}
