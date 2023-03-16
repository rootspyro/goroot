package goroot

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ROUTER

type Router struct {
	rules map[string]map[string]http.HandlerFunc
}

func(router *Router)findHandler(path string, method string) (http.HandlerFunc, bool, bool) {

	_, pathExists := router.rules[path]

	handler, methodExists := router.rules[path][method]

	return handler, pathExists, methodExists 
}

func(router *Router)ServeHTTP(w http.ResponseWriter, r *http.Request) {

	reqPath := r.URL.Path

	// global request log
	log.Printf("%s:%s - %s", r.Method, reqPath, r.Host)

	// Search if the path exists
	handler, pathExists, methodExists := router.findHandler(reqPath, r.Method)	

	// If path don't exists returns 404 not found
	if !pathExists {

		w.WriteHeader(http.StatusNotFound)
		return
	
	} 

	// If path exists but not the method then returns 405 method not allowed
	if !methodExists {

		w.WriteHeader(http.StatusMethodNotAllowed)
		return	
	}
	
	handler(w,r)
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
			rules: make(map[string]map[string]http.HandlerFunc),
		},
	}
}

func(s *Server)Handle(method, path string, handler http.HandlerFunc) {
	s.router.rules[path][method] = handler
}

func(s *Server)AddMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
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


