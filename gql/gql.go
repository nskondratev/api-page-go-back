package gql

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

type GraphQLHub struct {
	schema  graphql.Schema
	types   []graphql.Type
	queries graphql.Fields
}

func NewGraphQLHub() *GraphQLHub {
	return &GraphQLHub{
		types:   make([]graphql.Type, 0),
		queries: make(graphql.Fields),
	}
}

func (h *GraphQLHub) AddType(t graphql.Type) *GraphQLHub {
	h.types = append(h.types, t)
	return h
}

func (h *GraphQLHub) AddQuery(key string, q *graphql.Field) error {
	if _, ok := h.queries[key]; ok {
		return fmt.Errorf("query with key %s already exists", key)
	}
	h.queries[key] = q
	return nil
}

func (h *GraphQLHub) Execute(query string) (*graphql.Result, error) {
	result := graphql.Do(graphql.Params{
		Schema:        h.schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		return result, fmt.Errorf("graphql: fail to execute query: %v", result.Errors)
	}
	return result, nil
}

func (h *GraphQLHub) Compile() error {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: h.queries,
	})

	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		Types: h.types,
	})

	if err != nil {
		return err
	}
	h.schema = s
	return nil
}
