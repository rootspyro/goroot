package goroot

import "net/http"

type Server struct {
	port string
}

func Default() *Server{
	return &Server{
		port: "3000",
	}	
}

func New( _port string ) *Server {
	return &Server{
		port: _port,
	}
}

func(s Server) Listen() {
	http.ListenAndServe(":" + s.port, nil)
} 
