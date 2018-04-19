package main

import (
	"database/sql"
	"log"

	"github.com/davecgh/go-spew/spew"

	_ "github.com/go-sql-driver/mysql"
	"github.com/suppayami/goql/schema"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@/employees")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema, err := schema.BuildSQLSchema(db, schema.GetBuilder("mysql"), "employees")
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(schema)
}
