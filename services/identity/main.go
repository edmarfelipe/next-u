package main

import (
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/http"
)

func main() {
	ct := infra.NewContainer()

	server := http.NewServer(ct)
	server.Listen()
}
