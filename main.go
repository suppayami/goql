package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/handler"
	"github.com/suppayami/goql/resolver"
	"github.com/suppayami/goql/schema"
)

func main() {
	db, err := sql.Open("mysql", "root:12345678@/sakila")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlSchema, err := schema.BuildSQLSchema(db, schema.GetBuilder("mysql"))
	if err != nil {
		log.Fatal(err)
	}
	graphqlSchema, err := schema.SQLToGraphqlSchema(sqlSchema)
	if err != nil {
		log.Fatal(err)
	}

	exportGraphQL := flag.Bool("e", false, "Export graphql?")
	serveGraphQL := flag.Bool("s", false, "Serve graphql?")

	flag.Parse()

	if *exportGraphQL {
		fmt.Println(graphqlSchema)
	}

	if *serveGraphQL {
		schema, err := resolver.BuildSchema(sqlSchema, graphqlSchema)
		if err != nil {
			log.Fatal(err)
		}
		h := handler.New(&handler.Config{
			Schema:   schema,
			Pretty:   true,
			GraphiQL: true,
		})
		// fmt.Println(graphqlSchema)
		http.Handle("/", h)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	if !*exportGraphQL && !*serveGraphQL {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
