package main

import (
	"fmt"
	"github.com/4ttenji/tinisvr/handler"
	"net/http"
)

func main() {
	// set the http request handler using http.HandleFunc
	http.HandleFunc("/", handler.IndexHandler)

	fmt.Println("The server is running ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
