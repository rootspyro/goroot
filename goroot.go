package goroot

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

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
			rules: make(map[string]map[string]Handler),
		},
	}
}


// This function creates a new Endpoint in the router
func(s *Server)Endpoint(method, path string, handler Handler) {

	_, exists := s.router.rules[path]

	if !exists {
		s.router.rules[path] = make(map[string]Handler)
	}
	
	s.router.rules[path][method] = handler

}

// CRUD FUNCTIONS

// Creates a new endpoint in the rounter with the GET method
func(s *Server)Get(path string, handler Handler) {
	s.Endpoint("GET", path, handler)
}

// Creates a new endpoint in the rounter with the POST method
func(s *Server)Post(path string, handler Handler) {
	s.Endpoint("POST", path, handler)
}

// Creates a new endpoint in the rounter with the PUT method
func(s *Server)Put(path string, handler Handler) {
	s.Endpoint("GET", path, handler)
}

// Creates a new endpoint in the rounter with the PATCH method
func(s *Server)Patch(path string, handler Handler) {
	s.Endpoint("GET", path, handler)
}

// Creates a new endpoint in the rounter with the DELETE method
func(s *Server)Delete(path string, handler Handler) {
	s.Endpoint("GET", path, handler)
}


func(s *Server)AddMiddleware(handler Handler, middlewares ...Middleware) Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}

	return handler
}

func(s Server) Listen() {
	fmt.Println("Server listening in 0.0.0.0:" + s.port + "!")

	http.Handle("/", s.router)
	http.ListenAndServe(":" + s.port, nil)
} 


