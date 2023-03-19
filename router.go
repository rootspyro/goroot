package goroot

import (
	"net/http"
	"log"
)

type Router struct {
	rules map[string]map[string]Handler
}

func(router *Router)findHandler(path string, method string) (Handler, bool, bool) {

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
	
	handler(&Root{ writter: w, request: r })
}

type Route struct { 
	
}
