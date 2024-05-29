package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/wetee-dao/indexer/graph"
	"github.com/wetee-dao/indexer/store"
	"github.com/wetee-dao/indexer/util"
)

const defaultPort = "8881"

// 启动GraphQL服务器
// StartServer starts the GraphQL server.
func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Handle("/", playground.Handler("Wetee-Cache", "/gql"))
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers:  &graph.Resolver{},
		Directives: graph.NewDirectiveRoot(),
	}))
	router.Handle("/gql", srv)

	if util.IsFileExists(util.WORK_DIR+"/ser.pem") && util.IsFileExists(util.WORK_DIR+"/ser.key") {
		log.Printf("connect to https://localhost:%s/ for GraphQL playground", port)
		http.ListenAndServeTLS(":"+port, util.WORK_DIR+"/ser.pem", util.WORK_DIR+"/ser.key", router)
	} else {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		http.ListenAndServe(":"+defaultPort, router)
	}
}

func main() {
	// 初始化数据库
	err := store.DBInit(util.WORK_DIR + "/db")
	if err != nil {
		fmt.Println(err, "unable to start database")
		os.Exit(1)
	}
	StartServer()
}
