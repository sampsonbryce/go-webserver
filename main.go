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

type Person struct {
	Id   int    `json"id"`
	Name string `json:"name"`
}

func main() {
	options := server.ServerOptions{Host: connHost, Port: connPort}
	s := server.CreateServer(&options)

	s.HandleFunc("/bacon", func(server.Request) server.Response {
		p := Person{Id: 1, Name: "Dave"}
		return server.CreateResponse().SetStatus(200).SetJson(p)
	})

	fmt.Println("Starting server on " + connHost + ":" + connPort)

	log.Fatal(s.Listen())
}
