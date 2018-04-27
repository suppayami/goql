package resolver

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/suppayami/goql/schema"
)

// BuildSchema builds GraphQL handler & resolver
func BuildSchema(sqlSchema schema.SQLSchemaStruct, graphqlSchema schema.GraphqlSchema) (*graphql.Schema, error) {
	objectTypes := buildObjectTypes(sqlSchema, graphqlSchema)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: buildQueryType(sqlSchema, graphqlSchema, objectTypes),
	})
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

func buildObjectTypes(sqlSchema schema.SQLSchemaStruct, graphqlSchema schema.GraphqlSchema) map[string]*graphql.Object {
	objectTypes := make(map[string]*graphql.Object)
	// init object types
	for _, gql := range graphqlSchema.ObjectTypes {
		objectTypes[gql.Name] = graphql.NewObject(graphql.ObjectConfig{
			Name:   gql.Name,
			Fields: graphql.Fields{},
		})
	}
	// setup fields
	for _, gql := range graphqlSchema.ObjectTypes {
		objectType := objectTypes[gql.Name]
		for _, field := range gql.Fields {
			f := field
			objectType.AddFieldConfig(f.Name, &graphql.Field{
				Type: getGraphqlType(f, objectTypes),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if obj, ok := p.Source.(map[string]string); ok == true {
						if f.Type != schema.ObjectType {
							return obj[f.Name], nil
						}
						return nil, nil
					}
					return nil, nil
				},
			})
		}
	}
	return objectTypes
}

func buildQueryType(
	sqlSchema schema.SQLSchemaStruct,
	graphqlSchema schema.GraphqlSchema,
	objectTypes map[string]*graphql.Object,
) *graphql.Object {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   graphqlSchema.QueryType.Name,
		Fields: graphql.Fields{},
	})
	for _, queryField := range graphqlSchema.QueryType.Fields {
		qf := queryField
		rootQuery.AddFieldConfig(qf.Name, &graphql.Field{
			Type: getGraphqlType(qf, objectTypes),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				objectType := getGraphqlObjectType(graphqlSchema, qf.ObjectType)
				if len(qf.Arguments) == 0 {
					return objectType.Reader(make(map[string]interface{})), nil
				}
				return nil, nil
			},
		})
	}
	return rootQuery
}

func getGraphqlType(gql schema.GraphqlField, objectTypes map[string]*graphql.Object) graphql.Output {
	var gqlType graphql.Type
	switch gql.Type {
	case schema.ScalarID:
		gqlType = graphql.ID
	case schema.ScalarInt:
		gqlType = graphql.Int
	case schema.ScalarFloat:
		gqlType = graphql.Float
	case schema.ScalarBoolean:
		gqlType = graphql.Boolean
	case schema.ScalarString:
		gqlType = graphql.String
	case schema.ObjectType:
		gqlType = objectTypes[gql.ObjectType]
	}
	if !gql.Nullable {
		gqlType = graphql.NewNonNull(gqlType)
	}
	if gql.IsArray {
		gqlType = graphql.NewList(gqlType)
	}
	return gqlType
}

func getGraphqlObjectType(graphqlSchema schema.GraphqlSchema, objectType string) schema.GraphqlObjectType {
	for _, gql := range graphqlSchema.ObjectTypes {
		if strings.EqualFold(gql.Name, objectType) {
			return gql
		}
	}
	panic(fmt.Sprintf("ObjectType %s is missing", objectType))
}
