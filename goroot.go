package goroot

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rootspyro/goroot/cors"
	"github.com/rootspyro/goroot/pages"
)

// SERVER

type Server struct {
	port string
	router *Router
	config *Config
}

type Config struct {

	// Cors Config
	Cors *cors.Cors

	// List of Global Middlewares
	Middlewares []Middleware

	// Html rendering config
	Pages *pages.Pages
}

func New(config Config) *Server {

	var defPort string = "3000"

	// If the env PORT variable exists, use it as default port
	val, exists := os.LookupEnv("PORT")

	if exists {
		defPort = val
	}

	// If the -p flag is used, the port will be modified even if the env PORT variable exists.
	p := flag.String("p", defPort, "Server port")
	flag.Parse()

	if config.Cors == nil {
		config.Cors = cors.New(cors.Config{
			Origins: []string{"*"},
			Methods: []string{"GET", "POST", "PUT", "PATCH","DELETE"},
		})
	}

	return &Server{
		port: *p,
		config: &config,
		router: &Router{
			cors: config.Cors,
			middlewares: &config.Middlewares,
			pages: config.Pages,
			node: &Node{
				path: "/",
				actions: make(map[string]*Handler),
				children: make(map[string]*Node),
			},
		},
	}
}


// This function creates a new Endpoint in the router
func(s *Server)Endpoint(method, path string, handler Handler) {
	
	currentNode := s.router.node

	if path == "/" {
		currentNode.actions[method] = &handler
		return
	}

	ep := s.router.explodePath(path)

	for i, label := range ep {
		nextNode, exists := currentNode.children[label]	

		// If child node exists then update the current node for the next loop
		if exists {
			currentNode = nextNode
		}

		// If child node don't exists then create it
		if !exists { 
			currentNode.children[label] = &Node{
				path: label,
				actions: make(map[string]*Handler),
				children: make(map[string]*Node),
			}

			currentNode = currentNode.children[label]

		}

		// If is the last loop, then create the action in the current node
		if i == len(ep) - 1 {
			currentNode.path = label
			currentNode.actions[method] = &handler
			break
		}

	}

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

// Set the list of Global Middlewares
func(s Server) Middlewares(middlewares ...Middleware) {

	for _, midd := range middlewares {
		s.config.Middlewares = append(s.config.Middlewares, midd)
	}

}

