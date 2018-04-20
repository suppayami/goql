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
	Arguments  []GraphqlArgument
}

func (gql GraphqlField) String() string {
	var keywordFieldType, keywordArguments string

	if gql.Type == ObjectType {
		keywordFieldType = gql.ObjectType
	} else {
		keywordFieldType = string(gql.Type)
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
		return fmt.Sprintf(KeywordFieldDefaultValue, gql.Name, keywordFieldType, gql.DefaultValue)
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

// ConvertTypeSQLToGraphql converts SQL type to Graphql type
// TODO: different db has different types
func ConvertTypeSQLToGraphql(sqlType string) GraphqlType {
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

// IsIDField checks if given SQL field is an ID field
func IsIDField(sqlField SQLFieldStruct) bool {
	lowerField := strings.ToLower(sqlField.Field)
	if strings.HasSuffix(sqlField.Field, "ID") ||
		strings.HasSuffix(lowerField, "_id") {
		return true
	}
	return false
}

// ConvertFieldSQLtoGraphql converts an SQL field into a Graphql field
func ConvertFieldSQLtoGraphql(sqlField SQLFieldStruct) error {
	return nil
}
