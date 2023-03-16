package goroot

import (
	"flag"
	"net/http"
	"os"
)

type Server struct {
	port string
}

func New() *Server {

	var defPort string = "3000"

	val, exists := os.LookupEnv("PORT")

	if exists {
		defPort = val
	}

	p := flag.String("p", defPort, "Server port")
	flag.Parse()

	return &Server{
		port: *p,
	}
}

func(s Server) Listen() {
	http.ListenAndServe(":" + s.port, nil)
} 


