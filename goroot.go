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
	rules map[string]http.HandlerFunc
}

func(router *Router)findHandler(path string) (http.HandlerFunc, bool) {
	handler, exits := router.rules[path]
	return handler, exits
}

func(router *Router)ServeHTTP(w http.ResponseWriter, r *http.Request) {

	reqPath := r.URL.Path

	// Search if the path exists
	handler, exits := router.findHandler(reqPath)	

	// If path don't exists return 404 not found
	if !exits {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(reqPath + " not Found!"))
		return
	}

	// global request log
	log.Printf("%s:%s - %s", r.Method, r.URL.Path, r.Host)
	
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
			rules: make(map[string]http.HandlerFunc),
		},
	}
}

func(s *Server)Handle(path string, handler http.HandlerFunc) {
	s.router.rules[path] = handler
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


