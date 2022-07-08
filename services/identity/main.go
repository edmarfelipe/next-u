package main

import (
	"log"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/http"
)

func main() {
	config, err := infra.NewConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	server := http.SetupHTTPServer()
	server.Listen(config.ServerAddress())
}
