package server

import (
	"errors"
	"fmt"
	"net/http"
)

var p *pilot

type pilot struct {
	patterns map[string]http.Handler
}

func (p *pilot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := p.patterns[r.URL.Path]; ok {
		p.patterns[r.URL.Path].ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("did not match any handler"))
	}
}

func (p *pilot) add(pattern string, handler http.Handler) error {
	if pattern == "" {
		return errors.New("url pattern is not set")
	}

	if _, ok := p.patterns[pattern]; ok {
		return errors.New("duplicated url pattern")
	} else {
		p.patterns[pattern] = handler
		return nil
	}
}

func New() error {

	if p != nil {
		return errors.New("duplicated init process")
	} else {
		p = &pilot{
			patterns: make(map[string]http.Handler),
		}
		return nil
	}
}

func Add(pattern string, handler http.Handler) {
	p.add(pattern, handler)
}

func Run(addr string) {
	if p == nil {
		fmt.Println("Pilot is not existed. Initiation now.")
		p = &pilot{}
	}
	if err := http.ListenAndServe(addr, p); err != nil {
		fmt.Println(err)
	}
}
