package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/YukiOnishi1129/go-boilerplate-docker-graphql-postgres/app/util/initializer"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

const containerPort = "3000"

func main() {
	router := chi.NewRouter()

	srv, err := initializer.Init(router)
	if err != nil {
		panic(err)
	}

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	httpErr := http.ListenAndServe(":3000", router)
	if httpErr != nil {
		panic(httpErr)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", "4000")
	log.Fatal(http.ListenAndServe(":"+containerPort, nil))

}
