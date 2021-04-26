package server3

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var p *pilot

type pilot struct {
	patterns map[*Node]Handler
	nodeRoot *Node
}

type Context struct {
	w      http.ResponseWriter
	r      *http.Request
	Path   string
	Params map[string]string
}

func (c *Context) HTML(statusCode int, text string) {
	c.w.WriteHeader(statusCode)
	c.w.Write([]byte(text))
}

func (c *Context) Query(key string) string {
	return c.r.FormValue(key)
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

func (c *Context) Param(key string) string{
	return c.Params[key]
}

type Node struct {
	path     string
	part     string
	children map[string]*Node
	isWild   bool
}

func (n *Node) addChild(part string) *Node {
	c := &Node{
		part:     part,
		isWild:   strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*"),
		children: make(map[string]*Node),
	}
	n.children[part] = c
	return c
}

func (n *Node) insert(pattern string, parts []string, idx int) *Node {

	if idx == len(parts) {
		n.path = pattern
		return n
	}

	part := parts[idx]

	if _, ok := n.children[part]; !ok {
		return n.addChild(part).insert(pattern, parts, idx+1)
	}
	return n.children[part].insert(pattern, parts, idx+1)

}

// 这个写法建立在node 下不能同时存在通配符节点和普通节点的前提下，不然结果无法预期
func (n *Node) isMatch(part string) []string {

	res := make([]string, 0)
	for k, _ := range n.children {
		if k == part || strings.HasPrefix(k, ":") || strings.HasPrefix(k, "*") {
			res = append(res, k)
		}
	}

	return res
}

func (n *Node) search(parts []string, idx int) *Node {

	if len(parts) == idx {
		if n.path != "" {
			return n
		} else {
			return nil
		}
	}

	if len(n.isMatch(parts[idx])) != 0 {
		for _, v := range n.isMatch(parts[idx]) {
			if strings.HasPrefix(v, "*") {
				return n.children[v]
			}
			if dn := n.children[v].search(parts, idx+1); dn != nil {
				return dn
			}
		}
	}
	return nil

}

type H map[string]string

type Handler func(ctx *Context)

func (p *pilot) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := &Context{
		w,
		r,
		r.URL.Path,
		make(map[string]string),
	}

	raw_parts := strings.Split(c.Path, "/")
	part := make([]string, 0)
	for _, v := range raw_parts {
		if v != "" {
			part = append(part, v)
		}
	}

	n := p.nodeRoot.search(part, 0)

	rawPaths := strings.Split(n.path, "/")
	rawMatch := strings.Split(c.Path, "/")
	for k := range rawPaths {
		if strings.HasPrefix(rawPaths[k], ":") {
			c.Params[rawPaths[k][1:]] = rawMatch[k]
		} else if strings.HasPrefix(rawPaths[k], "*") {
			c.Params[rawPaths[k][1:]] = strings.Join(rawMatch[k:], "/")
		}
	}

	if n != nil {
		p.patterns[n](c)
	} else {
		c.w.WriteHeader(http.StatusNotFound)
		c.w.Write([]byte("don't match any url"))
	}
}

func GET(pattern string, handler func(c *Context)) error {
	// process the pattern string
	if !strings.HasPrefix(pattern, "/") {
		return errors.New("invalid pattern")
	}

	raw_parts := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, part := range raw_parts {
		if part != "" {
			parts = append(parts, part)
		}
	}

	n := p.nodeRoot.insert(pattern, parts, 0)
	p.patterns[n] = handler
	return nil
}

func New() {
	p = &pilot{
		patterns: make(map[*Node]Handler),
		nodeRoot: &Node{
			path:     "",
			part:     "",
			children: make(map[string]*Node),
			isWild:   false,
		},
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
