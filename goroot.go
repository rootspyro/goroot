package goroot

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

// ROUTER

type Router struct {
	rules map[string]http.HandlerFunc
}

func(router *Router)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}

// SERVER

type Server struct {
	port string
	router *Router
}


func New() *Server {

	var defPort string = "3000"

	// If the env PORT variable exists, use it as default port
	val, exists := os.LookupEnv("PORT")

	if exists {
		defPort = val
	}

	// If the -p flag is used, the port will be modified even if the env PORT variable exists.
	p := flag.String("p", defPort, "Server port")
	flag.Parse()

	return &Server{
		port: *p,
		router: &Router{
			rules: make(map[string]http.HandlerFunc),
		},
	}
}

func(s Server) Listen() {
	fmt.Println("Server listening in 0.0.0.0:" + s.port + "!")

	http.Handle("/", s.router)
	http.ListenAndServe(":" + s.port, nil)
} 


