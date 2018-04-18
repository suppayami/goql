package schema

import (
	"strings"
)

// ScalarType is a scalar type in GraphQL
type ScalarType string

// Scalar Types in GraphQL
const (
	ScalarInt     ScalarType = "Int"
	ScalarFloat   ScalarType = "Float"
	ScalarString  ScalarType = "String"
	ScalarBoolean ScalarType = "Boolean"
	ScalarID      ScalarType = "ID"
)

// GraphqlSchemaBuilder pipes DBSchema into a barebone GraphqlSchema.
type GraphqlSchemaBuilder interface{}

// ConvertTypeSQLToGraphql converts SQL type to Graphql type
func ConvertTypeSQLToGraphql(sqlType string) ScalarType {
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
	if len(sqlField.Key) > 0 {
		return true
	}
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
