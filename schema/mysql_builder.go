package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

// MySQLSchemaBuilder implements SQLSchemaBuilder
type MySQLSchemaBuilder struct{}

type mySQLField struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

// QueryTables implementation
func (builder MySQLSchemaBuilder) QueryTables(db *sql.DB) ([]SQLTableStruct, error) {
	tables := []SQLTableStruct{}
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return tables, err
	}
	defer rows.Close()
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return tables, err
		}
		table := SQLTableStruct{
			Name:   tableName,
			Fields: []SQLFieldStruct{},
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return tables, err
	}
	return tables, nil
}

// QueryFields implementation
func (builder MySQLSchemaBuilder) QueryFields(db *sql.DB, tableName string) ([]SQLFieldStruct, error) {
	fields := []SQLFieldStruct{}
	rows, err := db.Query(fmt.Sprintf("DESCRIBE %s", tableName))
	if err != nil {
		return fields, err
	}
	defer rows.Close()
	for rows.Next() {
		var fieldStruct mySQLField
		if err := rows.Scan(
			&fieldStruct.Field,
			&fieldStruct.Type,
			&fieldStruct.Null,
			&fieldStruct.Key,
			&fieldStruct.Default,
			&fieldStruct.Extra,
		); err != nil {
			fmt.Println(rows)
			return fields, err
		}
		field := SQLFieldStruct{
			Field: fieldStruct.Field,
			Null:  strings.EqualFold(fieldStruct.Null, "yes"),
			Type:  fieldStruct.Type,
		}
		fields = append(fields, field)
	}
	if err := rows.Err(); err != nil {
		return fields, err
	}
	return fields, nil
}
