package schema

// SQLSchemaBuilder queries the database schema and builds it into a readable struct,
// allows GraphqlSchemaBuilder build a barebone Graphql schema from database structure.
type SQLSchemaBuilder interface{}

// SQLFieldStruct describes a field in table of database.
type SQLFieldStruct struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
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
