package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

// SQLSchemaBuilder queries the database schema and builds it into a readable struct,
// allows GraphqlSchemaBuilder build a barebone Graphql schema from database structure.
type SQLSchemaBuilder interface {
	// QueryTables should only returns a slice of SQLTableStruct without the Fields.
	// The fields will be appended in the main builder function.
	QueryTables(db *sql.DB) ([]*SQLTableStruct, error)

	// QueryFields should map the table description from database to a slice of
	// SQLFieldStruct.
	QueryFields(db *sql.DB, tableName string) ([]*SQLFieldStruct, error)
}

// SQLFieldStruct describes a field in table of database.
type SQLFieldStruct struct {
	Field        string
	Type         string
	Null         bool
	IsPrimaryKey bool
	IsForeignKey bool
}

// SQLTableStruct describes a table in database.
type SQLTableStruct struct {
	Name          string
	Fields        []*SQLFieldStruct
	Relationships []*SQLRelationshipStruct
	IsManyToMany  bool
}

// SQLRelationshipStruct describes a relationship between tables.
type SQLRelationshipStruct struct {
	Table      *SQLTableStruct
	ForeignKey string
	LocalKey   string
	Null       bool
	HasMany    bool
}

// SQLSchemaStruct describes database schema.
type SQLSchemaStruct struct {
	Tables []*SQLTableStruct
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
func BuildSQLSchema(db *sql.DB, builder SQLSchemaBuilder) (SQLSchemaStruct, error) {
	schema := SQLSchemaStruct{
		Tables: []*SQLTableStruct{},
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
		setupRelationships(tables, table)
		table.IsManyToMany = isManyToManyTable(table)
		schema.Tables = append(schema.Tables, table)
	}
	return schema, nil
}

func setupRelationships(tableList []*SQLTableStruct, table *SQLTableStruct) {
	for i := range table.Fields {
		field := table.Fields[i]
		if !IsKey(*field) {
			continue
		}
		modelName := TableName(field.Field)
		if strings.EqualFold(table.Name, modelName) {
			continue
		}
		for j := range tableList {
			foundTable := tableList[j]
			if !strings.EqualFold(foundTable.Name, modelName) {
				continue
			}
			relForeign := SQLRelationshipStruct{
				Table:      foundTable,
				ForeignKey: field.Field,
				LocalKey:   PrimaryKey(foundTable.Name),
				Null:       field.Null,
			}
			table.Relationships = append(table.Relationships, &relForeign)
			relLocal := SQLRelationshipStruct{
				Table:      table,
				ForeignKey: PrimaryKey(foundTable.Name),
				LocalKey:   field.Field,
				HasMany:    true,
				Null:       true,
			}
			foundTable.Relationships = append(foundTable.Relationships, &relLocal)
			field.IsForeignKey = true
		}
	}
}

// TODO: Many-to-many relationship should be checked by some conventions
func isManyToManyTable(table *SQLTableStruct) bool {
	if len(table.Relationships) < 2 {
		return false
	}
	for _, field := range table.Fields {
		if strings.EqualFold(field.Field, fmt.Sprintf("%s_id", table.Name)) {
			return false
		}
	}
	return true
}
