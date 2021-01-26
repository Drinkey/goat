package main

import (
	"log"
	"net/http"

	"github.com/Drinkey/goat/routers"
)

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
