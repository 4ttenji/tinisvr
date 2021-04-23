package handler

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	response := "It's index page."
	w.Header().Set("accepted", "true")
	w.WriteHeader(200)
	w.Write([]byte(response))

}
