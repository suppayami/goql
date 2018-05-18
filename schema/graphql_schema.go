package schema

import (
	"fmt"
	"strings"
)

// GraphqlType is a type in GraphQL
type GraphqlType string

// Scalar Types in GraphQL
const (
	ScalarInt     GraphqlType = "Int"
	ScalarFloat   GraphqlType = "Float"
	ScalarString  GraphqlType = "String"
	ScalarBoolean GraphqlType = "Boolean"
	ScalarID      GraphqlType = "ID"
	ObjectType    GraphqlType = "ObjectType"
)

// Graphql Keyword
const (
	KeywordType     string = "type"
	KeywordInput    string = "input"
	KeywordSchema   string = "schema"
	KeywordQuery    string = "Query"
	KeywordMutation string = "Mutation"

	KeywordArray           string = "[%s]"
	KeywordNonNullableType string = "%s!"

	KeywordField             string = "%s: %s"
	KeywordFieldArguments    string = "%s(%s): %s"
	KeywordFieldDefaultValue string = "%s: %s = %s"
)

// GraphqlSchemaBuilder pipes DBSchema into a barebone GraphqlSchema.
type GraphqlSchemaBuilder interface{}

// GraphqlField indicates a field.
// ObjectType is an optional field, needed when Type is ObjectType.
type GraphqlField struct {
	Name       string
	Type       GraphqlType
	ObjectType string
	Nullable   bool
	IsArray    bool
	Arguments  []GraphqlArgument
}

func (gql GraphqlField) String() string {
	var keywordFieldType, keywordArguments string

	if gql.Type == ObjectType {
		keywordFieldType = gql.ObjectType
	} else {
		keywordFieldType = string(gql.Type)
	}

	if gql.IsArray {
		keywordFieldType = fmt.Sprintf(KeywordArray, keywordFieldType)
	}

	if !gql.Nullable {
		keywordFieldType = fmt.Sprintf(KeywordNonNullableType, keywordFieldType)
	}

	if len(gql.Arguments) > 0 {
		args := make([]string, 0, len(gql.Arguments))
		for _, arg := range gql.Arguments {
			args = append(args, arg.String())
		}
		keywordArguments = strings.Join(args, ", ")
		return fmt.Sprintf(KeywordFieldArguments, gql.Name, keywordArguments, keywordFieldType)
	}

	return fmt.Sprintf(KeywordField, gql.Name, keywordFieldType)
}

// GraphqlArgument is used for a field's argument.
// DefaultValue is always a string and casted based on Type.
type GraphqlArgument struct {
	Name         string
	Type         GraphqlType
	ObjectType   string
	Nullable     bool
	DefaultValue string
}

func (gql GraphqlArgument) String() string {
	var keywordFieldType string

	if gql.Type == ObjectType {
		keywordFieldType = gql.ObjectType
	} else {
		keywordFieldType = string(gql.Type)
	}

	if len(gql.DefaultValue) > 0 {
		defaultValue := gql.DefaultValue
		if gql.Type == ScalarString {
			defaultValue = fmt.Sprintf("\"%s\"", defaultValue)
		}
		return fmt.Sprintf(KeywordFieldDefaultValue, gql.Name, keywordFieldType, defaultValue)
	}

	if !gql.Nullable {
		keywordFieldType = fmt.Sprintf(KeywordNonNullableType, keywordFieldType)
	}

	return fmt.Sprintf(KeywordField, gql.Name, keywordFieldType)
}

// GraphqlObjectType describes an object type in Graphql.
type GraphqlObjectType struct {
	Name   string
	Fields []GraphqlField
}

func (gql GraphqlObjectType) String() string {
	fields := make([]string, 0, len(gql.Fields))
	for _, field := range gql.Fields {
		fields = append(fields, fmt.Sprintf("\t%s", field.String()))
	}
	return fmt.Sprintf("%s %s {\n%s\n}", KeywordType, gql.Name, strings.Join(fields, "\n"))
}

// GraphqlSchema describes Graphql schema
type GraphqlSchema struct {
	QueryType    GraphqlObjectType
	MutationType GraphqlObjectType
	ObjectTypes  []GraphqlObjectType
}

func (gql GraphqlSchema) String() string {
	objectTypes := make([]string, 0, len(gql.ObjectTypes))
	objectTypes = append(objectTypes, gql.QueryType.String())
	objectTypes = append(objectTypes, gql.MutationType.String())
	for _, objectType := range gql.ObjectTypes {
		objectTypes = append(objectTypes, objectType.String())
	}
	schemaTxt := fmt.Sprintf("%s {\n\tquery: Query\n\tmutation: Mutation\n}\n\n", KeywordSchema)
	return fmt.Sprintf("%s%s", schemaTxt, strings.Join(objectTypes, "\n\n"))
}

// SQLToGraphqlSchema converts SQL Schema to Graphql Schema
func SQLToGraphqlSchema(sqlSchema SQLSchemaStruct) (GraphqlSchema, error) {
	schema := GraphqlSchema{
		QueryType: GraphqlObjectType{
			Name:   "Query",
			Fields: []GraphqlField{},
		},
		MutationType: GraphqlObjectType{
			Name:   "Mutation",
			Fields: []GraphqlField{},
		},
		ObjectTypes: []GraphqlObjectType{},
	}

	for _, sqlTable := range sqlSchema.Tables {
		// if sqlTable.IsManyToMany {
		// 	continue
		// }
		objectType := sqlToGraphqlObjectType(sqlTable)
		queryFields := sqlToGraphqlQueryFields(sqlTable)
		schema.ObjectTypes = append(schema.ObjectTypes, objectType)
		for _, queryField := range queryFields {
			schema.QueryType.Fields = append(schema.QueryType.Fields, queryField)
		}
	}

	return schema, nil
}

// ConvertTypeSQLToGraphql converts SQL type to Graphql type
// TODO: different db has different types
func sqlToGraphqlType(sqlType string) GraphqlType {
	typeLower := strings.ToLower(sqlType)
	if strings.Contains(typeLower, "int") {
		return ScalarInt
	}
	if strings.Contains(typeLower, "float") ||
		strings.Contains(typeLower, "double") {
		return ScalarFloat
	}
	if strings.Contains(typeLower, "bit") {
		return ScalarBoolean
	}
	return ScalarString
}

func sqlToGraphqlObjectType(sqlTable *SQLTableStruct) GraphqlObjectType {
	objectType := GraphqlObjectType{
		Name:   SQLToGraphqlObjectName(sqlTable.Name),
		Fields: []GraphqlField{},
	}

	for _, sqlField := range sqlTable.Fields {
		gqlType := sqlToGraphqlType(sqlField.Type)
		if sqlField.IsForeignKey {
			continue
		}
		if sqlField.IsPrimaryKey {
			gqlType = ScalarID
		}
		field := GraphqlField{
			Name:     SQLToGraphqlFieldName(sqlField.Field),
			Type:     gqlType,
			IsArray:  false,
			Nullable: sqlField.Null,
		}
		objectType.Fields = append(objectType.Fields, field)
	}

	for _, sqlRelationship := range sqlTable.Relationships {
		if !sqlRelationship.Table.IsManyToMany {
			field := GraphqlField{
				Name:       SQLToGraphqlFieldName(sqlRelationship.Table.Name),
				Type:       ObjectType,
				ObjectType: SQLToGraphqlObjectName(sqlRelationship.Table.Name),
				IsArray:    sqlRelationship.HasMany,
				Nullable:   sqlRelationship.Null,
			}
			if field.IsArray {
				field.Name = ArrayFieldName(field.Name)
			}
			objectType.Fields = append(objectType.Fields, field)
			continue
		}
		// for _, manyToMany := range sqlRelationship.Table.Relationships {
		// 	if manyToMany.Table == sqlTable {
		// 		continue
		// 	}
		// 	field := GraphqlField{
		// 		Name:       ArrayFieldName(SQLToGraphqlFieldName(manyToMany.Table.Name)),
		// 		Type:       ObjectType,
		// 		ObjectType: SQLToGraphqlObjectName(manyToMany.Table.Name),
		// 		IsArray:    true,
		// 		Nullable:   true,
		// 	}
		// 	objectType.Fields = append(objectType.Fields, field)
		// }
	}

	return objectType
}

func sqlToGraphqlQueryFields(sqlTable *SQLTableStruct) []GraphqlField {
	queryFields := []GraphqlField{}
	queryFields = append(queryFields, GraphqlField{
		Name:       ArrayFieldName(SQLToGraphqlFieldName(sqlTable.Name)),
		Type:       ObjectType,
		ObjectType: SQLToGraphqlObjectName(sqlTable.Name),
		IsArray:    true,
		Nullable:   true,
		Arguments: []GraphqlArgument{
			GraphqlArgument{
				Name:         "first",
				Type:         ScalarInt,
				Nullable:     true,
				DefaultValue: "0",
			},

			GraphqlArgument{
				Name:         "offset",
				Type:         ScalarInt,
				Nullable:     true,
				DefaultValue: "5",
			},
		},
	})
	queryFields = append(queryFields, GraphqlField{
		Name:       SQLToGraphqlFieldName(sqlTable.Name),
		Type:       ObjectType,
		ObjectType: SQLToGraphqlObjectName(sqlTable.Name),
		IsArray:    false,
		Nullable:   true,
		Arguments: []GraphqlArgument{
			GraphqlArgument{
				Name:     PrimaryKey(sqlTable.Name),
				Type:     ScalarID,
				Nullable: false,
			},
		},
	})
	return queryFields
}
