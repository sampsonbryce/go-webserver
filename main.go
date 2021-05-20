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

	s.HandleFunc("/json", func(request server.Request) server.Response {
		if request.Method == "GET" {
			p := Person{Id: 1, Name: "Dave"}
			return server.CreateResponse().SetStatus(200).SetJson(p)
		} else if request.Method == "POST" {
			return server.CreateResponse().SetStatus(204)
		}

		return server.CreateResponse().SetStatus(404)
	})

	s.HandleFunc("/normal", func(request server.Request) server.Response {
		return server.CreateResponse().SetStatus(200).SetBody([]byte("helloworld"))
	})

	fmt.Println("Starting server on " + connHost + ":" + connPort)

	log.Fatal(s.Listen())
}
