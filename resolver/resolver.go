package resolver

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/suppayami/goql/schema"
)

// BuildSchema builds GraphQL handler & resolver
func BuildSchema(db *sql.DB, sqlSchema schema.SQLSchemaStruct, graphqlSchema schema.GraphqlSchema) (*graphql.Schema, error) {
	objectTypes := buildObjectTypes(db, sqlSchema, graphqlSchema)
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: buildQueryType(db, sqlSchema, graphqlSchema, objectTypes),
	})
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

func buildObjectTypes(
	db *sql.DB,
	sqlSchema schema.SQLSchemaStruct,
	graphqlSchema schema.GraphqlSchema,
) map[string]*graphql.Object {
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
		gqlName := gql.Name
		for _, field := range gql.Fields {
			f := field
			objectType.AddFieldConfig(f.Name, &graphql.Field{
				Type: getGraphqlType(f, objectTypes),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if obj, ok := p.Source.(map[string]string); ok == true {
						if f.Type != schema.ObjectType {
							return obj[f.Name], nil
						}

						// object type
						table := getSQLTable(sqlSchema, schema.GraphqlToSQLTableName(f.ObjectType))
						reader := makeReader(db, table)
						m := make(map[string]interface{})
						key := ""
						if f.IsArray {
							key = schema.PrimaryKey(schema.GraphqlToSQLTableName(gqlName))
						} else {
							key = schema.PrimaryKey(table.Name)
						}
						m[key] = obj[schema.SQLToGraphqlFieldName(key)]
						results := reader(m)
						if f.IsArray {
							return results, nil
						}
						return results[0], nil
					}
					return nil, nil
				},
			})
		}
	}
	return objectTypes
}

func buildQueryType(
	db *sql.DB,
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
		table := getSQLTable(sqlSchema, schema.GraphqlToSQLTableName(qf.ObjectType))
		reader := makeReader(db, table)
		rootQuery.AddFieldConfig(qf.Name, &graphql.Field{
			Type: getGraphqlType(qf, objectTypes),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return reader(make(map[string]interface{})), nil
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

func getSQLTable(sqlSchema schema.SQLSchemaStruct, tableName string) *schema.SQLTableStruct {
	for _, sqlTable := range sqlSchema.Tables {
		if strings.EqualFold(sqlTable.Name, tableName) {
			return sqlTable
		}
	}
	panic(fmt.Sprintf("Table %s is missing", tableName))
}
