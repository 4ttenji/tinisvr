package main

import (
	"fmt"
	"github.com/4ttenji/tinisvr/handler"
	"github.com/4ttenji/tinisvr/server2"
	"github.com/4ttenji/tinisvr/server3"
	"net/http"
)

func main1() {
	// set the http request handler using http.HandleFunc
	http.HandleFunc("/", handler.IndexHandler)

	fmt.Println("The server is running ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main2() {

	// init the internal http server
	_ = server2.New()

	_ = server2.GET("/", func(c *server2.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	_ = server2.GET("/hello", func(c *server2.Context) {
		// expect /hello?name=ten
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	_ = server2.POST("/login", func(c *server2.Context) {
		c.JSON(http.StatusOK, server2.H{
			"password": c.PostForm("password"),
			"username": c.PostForm("username"),
		})
	})

	server2.Run(":8085")

}

func main3() {

	server3.New()
	var err error
	if err = server3.GET("/", func(c *server3.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	});err != nil{
		fmt.Println(err)
	}

	if err = server3.GET("/hello/:name", func(c *server3.Context) {
		// expect /hello/ten
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	});err != nil{
		fmt.Println(err)
	}

	if err = server3.GET("/hello/:name", func(c *server3.Context) {
		// expect /hello/ten
		c.String(http.StatusOK, "hello %s, I'm at %s\n", c.Param("name"), c.Path)
	});err != nil{
		fmt.Println(err)
	}

	if err = server3.GET("/hello", func(c *server3.Context) {
		// expect /hello?name=ten
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	});err != nil{
		fmt.Println(err)
	}

	if err = server3.GET("/assets/*filepath", func(c *server3.Context) {
		c.JSON(http.StatusOK, server3.H{"filepath": c.Param("filepath")})
	});err != nil{
		fmt.Println(err)
	}

	server3.Run(":9998")

}

func main() {
	// main1()
	// main2()
	main3()
}
