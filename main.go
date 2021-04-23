package main

import (
	"fmt"
	"github.com/4ttenji/tinisvr/handler"
	"github.com/4ttenji/tinisvr/server"
	"net/http"
)

func main_old() {
	// set the http request handler using http.HandleFunc
	http.HandleFunc("/", handler.IndexHandler)

	fmt.Println("The server is running ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {

	// init the internel http server
	_ = server.New()

	server.Add("/startup", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Got it! The req to %s", r.URL.Path)
	}))

	server.Add("/quit", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Got it! The req to %s, bye-bye", r.URL.Path)
	}))

	server.Run(":8085")

}
