package server2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var p *pilot

type pilot struct {
	patterns     map[string]Handler
	postPatterns map[string]Handler
}

type Context struct {
	w    http.ResponseWriter
	r    *http.Request
	Path string
}

type H map[string]string

type Handler func(ctx *Context)

func (c *Context) PostForm(key string) string {
	return c.r.PostFormValue(key)
}

func (c *Context) Query(key string) string {
	return c.r.FormValue(key)
}

func (c *Context) HTML(statusCode int, text string) {
	c.w.WriteHeader(statusCode)
	c.w.Write([]byte(text))
}

func (c *Context) JSON(statusCode int, data H) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(statusCode)
	jByte, _ := json.Marshal(data)
	c.w.Write(jByte)
}

func (c *Context) String(statusCode int, text string, a ...interface{}) {
	fmt.Println(a)
	c.w.WriteHeader(statusCode)

	fmt.Fprintf(c.w, text, a...)
}

func (p *pilot) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := &Context{
		w,
		r,
		r.URL.Path,
	}
	switch c.r.Method {
	case "GET":
		if _, ok := p.patterns[r.URL.Path]; ok {
			p.patterns[r.URL.Path](c)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("did not match any handler"))
		}
	case "POST":
		if _, ok := p.postPatterns[r.URL.Path]; ok {
			p.postPatterns[r.URL.Path](c)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("did not match any handler"))
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}
}

func New() error {

	if p != nil {
		return errors.New("duplicated init process")
	} else {
		p = &pilot{
			patterns:     make(map[string]Handler),
			postPatterns: make(map[string]Handler),
		}
		return nil
	}
}

func GET(pattern string, handleFunc func(c *Context)) error {
	if pattern == "" {
		return errors.New("url pattern is not set")
	}

	if _, ok := p.patterns[pattern]; ok {
		return errors.New("duplicated url pattern")
	} else {
		p.patterns[pattern] = handleFunc
		return nil
	}
}

func POST(pattern string, handleFunc func(c *Context)) error {
	if pattern == "" {
		return errors.New("url pattern is not set")
	}

	if _, ok := p.postPatterns[pattern]; ok {
		return errors.New("duplicated url pattern")
	} else {
		p.postPatterns[pattern] = handleFunc
		return nil
	}
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
