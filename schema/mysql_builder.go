package schema

import "database/sql"

// MySQLSchemaBuilder implements SQLSchemaBuilder
type MySQLSchemaBuilder struct{}

// QueryTables implementation
func (builder MySQLSchemaBuilder) QueryTables(db *sql.DB) ([]SQLTableStruct, error) {
	return nil, nil
}

// QueryFields implementation
func (builder MySQLSchemaBuilder) QueryFields(db *sql.DB, tableName string) ([]SQLFieldStruct, error) {
	return nil, nil
}
