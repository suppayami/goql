package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/suppayami/goql/schema"
)

func main() {
	db, err := sql.Open("mysql", "root:12345678@/employees")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlSchema, err := schema.BuildSQLSchema(db, schema.GetBuilder("mysql"))
	if err != nil {
		log.Fatal(err)
	}
	graphqlSchema, err := schema.SQLToGraphqlSchema(sqlSchema)
	fmt.Println(graphqlSchema.String())
}
