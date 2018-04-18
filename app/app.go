package app

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go/relay"
)

// StartApp starts shit
func StartApp(db *sql.DB, addr string) {
	graphiql, err := getGraphiQLTemplate("./graphiql.html")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	// TODO: router to another file
	// TODO: put graphql schema here after build schema
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(graphiql)
	})
	mux.Handle("/graphql", &relay.Handler{})

	log.Fatal(http.ListenAndServe(addr, mux))
}

func getGraphiQLTemplate(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
