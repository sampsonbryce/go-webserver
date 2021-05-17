package main

import (
	"fmt"
	"log"

	"github.com/sampsonbryce/go-webserver/server"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

func main() {
	options := server.ServerOptions{Host: connHost, Port: connPort}
	s := server.CreateServer(&options)

	s.HandleFunc("/bacon", func(server.Request) server.Response {
		return server.Response{StatusCode: 200}
	})

	fmt.Println("Starting server on " + connHost + ":" + connPort)

	log.Fatal(s.Listen())
}
