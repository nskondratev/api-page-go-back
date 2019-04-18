package pages

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/nskondratev/api-page-go-back/gql"
)

var GraphQLType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Page",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"text": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

func RegisterGraphQLQueries(ps Store, hub *gql.GraphQLHub) error {
	pageByIdQuery := &graphql.Field{
		Type:        GraphQLType,
		Description: "Get page by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: getByIdResolver(ps),
	}
	if err := hub.AddQuery("page", pageByIdQuery); err != nil {
		return err
	}
	return nil
}

func getByIdResolver(ps Store) func(graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(int)
		if !ok {
			return nil, errors.New("graphql: cannot parse id argument")
		}
		page, err := ps.GetById(uint64(id))
		if err != nil {
			return nil, err
		}
		return page, nil
	}
}
