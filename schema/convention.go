package schema

import (
	"fmt"
	"strings"

	"github.com/hungneox/stringutils"

	"github.com/jinzhu/inflection"
)

const (
	sqlIDSuffix = "_id"
)

// IsPrimaryKey check if the field is primary key, used for relationship
func IsPrimaryKey(sqlTable SQLTableStruct, sqlField SQLFieldStruct) bool {
	return strings.EqualFold(fmt.Sprintf("%s%s", sqlTable.Name, sqlIDSuffix), sqlField.Field)
}

// IsForeignKey check if the field is foreign key
func IsForeignKey(sqlTable SQLTableStruct, sqlField SQLFieldStruct) bool {
	return IsKey(sqlField) && !IsPrimaryKey(sqlTable, sqlField)
}

// IsKey check if the field is a key
func IsKey(sqlField SQLFieldStruct) bool {
	return strings.HasSuffix(sqlField.Field, sqlIDSuffix)
}

// PrimaryKey returns primary key name for table
func PrimaryKey(tableName string) string {
	return fmt.Sprintf("%s%s", tableName, sqlIDSuffix)
}

// TableName returns table name from key name
func TableName(key string) string {
	return strings.TrimSuffix(key, sqlIDSuffix)
}

// SQLToGraphqlFieldName returns case for graphql field
func SQLToGraphqlFieldName(fieldName string) string {
	return stringutils.CamelCase(fieldName)
}

// GraphqlToSQLFieldName returns case for sql field
func GraphqlToSQLFieldName(fieldName string) string {
	return stringutils.SnakeCase(fieldName)
}

// SQLToGraphqlObjectName returns case for graphql object
func SQLToGraphqlObjectName(tableName string) string {
	return stringutils.PascalCase(tableName)
}

// GraphqlToSQLTableName returns case for sql table
func GraphqlToSQLTableName(objectName string) string {
	return stringutils.SnakeCase(objectName)
}

// ArrayFieldName returns name for an array field
func ArrayFieldName(fieldName string) string {
	return inflection.Plural(fieldName)
}
