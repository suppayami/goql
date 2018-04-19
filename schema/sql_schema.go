package schema

import (
	"database/sql"
	"fmt"
)

// SQLSchemaBuilder queries the database schema and builds it into a readable struct,
// allows GraphqlSchemaBuilder build a barebone Graphql schema from database structure.
type SQLSchemaBuilder interface {
	QueryTables(db *sql.DB) ([]SQLTableStruct, error)
	QueryFields(db *sql.DB, tableName string) ([]SQLFieldStruct, error)
}

// SQLFieldStruct describes a field in table of database.
type SQLFieldStruct struct {
	Field string
	Type  string
	Null  string
}

// SQLTableStruct describes a table in database.
type SQLTableStruct struct {
	Name   string
	Fields []SQLFieldStruct
}

// SQLSchemaStruct describes database schema.
type SQLSchemaStruct struct {
	Name   string
	Tables []SQLTableStruct
}

// GetBuilder switches SQL driver to a SQLSchemaBuilder
func GetBuilder(driver string) SQLSchemaBuilder {
	switch driver {
	case "mysql":
		return MySQLSchemaBuilder{}
	default:
		panic(fmt.Sprintf("%s driver is not supported", driver))
	}
}

// BuildSQLSchema builds a SQL Schema from given connecting database
func BuildSQLSchema(db *sql.DB, builder SQLSchemaBuilder, dbName string) (SQLSchemaStruct, error) {
	schema := SQLSchemaStruct{
		Name:   dbName,
		Tables: nil,
	}
	tables, err := builder.QueryTables(db)
	if err != nil {
		return schema, err
	}
	for _, table := range tables {
		fields, err := builder.QueryFields(db, table.Name)
		if err != nil {
			return schema, err
		}
		table.Fields = fields
		schema.Tables = append(schema.Tables, table)
	}
	return schema, nil
}
