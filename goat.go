package main

import (
	"log"
	"net/http"

	"github.com/Drinkey/goat/routers"
)

// @title GoAt API Document
// @version 1.0
// @description GoAt API Document. Please notice the responses described in this document is response "data" field's value, not the real response
// @host localhost:8090
// @BasePath /api/v1
// @query.collection.format multi
func main() {
	apiServerHandler := routers.InitRouter()
	apiServer := &http.Server{
		Addr:    ":8090",
		Handler: apiServerHandler,
	}
	err := apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
