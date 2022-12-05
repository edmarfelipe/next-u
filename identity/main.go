package main

import (
	"github.com/edmarfelipe/next-u/identity/infra"
	"github.com/edmarfelipe/next-u/identity/infra/http/server"
)

func main() {
	ct := infra.NewContainer()
	server := server.New(ct)
	server.Listen()
}
