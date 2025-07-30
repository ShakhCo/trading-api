package main

import (
	"fmt"
	"trading_api/db"
	routes "trading_api/handlers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Handler for /hi
func hiHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi from FastHTTP!")
}

func main() {

	db.InitDB()

	r := router.New()

	r.GET("/hi", hiHandler)
	r.POST("/register", routes.RegisterUserHandler)
	r.GET("/users", routes.GetAllUsersHandler)

	fmt.Println("FastHTTP server is running on :9000")
	if err := fasthttp.ListenAndServe(":9000", r.Handler); err != nil {
		panic("Error starting server: " + err.Error())
	}
}
