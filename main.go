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
	"github.com/suppayami/goql/env"
	"github.com/suppayami/goql/resolver"
	"github.com/suppayami/goql/schema"
)

func main() {
	var e env.Env

	e.ReadEnv()

	// TODO: Auto detect database server (mysql|postgres|mongodb?)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s", e.Username, e.Password, e.Host, e.Port, e.Database))

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
		schema, err := resolver.BuildSchema(db, sqlSchema, graphqlSchema)
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
