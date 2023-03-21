package goroot

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/rootspyro/goroot/cors"
)

type Router struct {
	//Base Node
	node *Node

	// Cors Configuration
	cors *cors.Cors

	// Global Middlewares
	middlewares *[]Middleware
}

type Node struct {
	path		 string
	actions	 map[string]*Handler
	children map[string]*Node
}

func(router *Router)findHandler(path, method string, root *Root) (Handler,  bool, bool) {
 
	currentNode := router.node

	if path != "/" { 

		for _, label := range router.explodePath(path) {

			nextNode, exists := currentNode.children[label]

			if !exists {

				if currentNode.path == label {

					break

				} else {

					// Boolean value for search if exists an children with a parameter as path. Example = /{userID}
					founded := false

					for path, node := range currentNode.children {

						if router.isParameter(path) {

							founded = true
							root.RequestParams[router.clearParam(path)]	= label 
							currentNode = node

							break

						}
						
					}

					if founded {
						// if node exists and has children then continue
						if len(currentNode.children) > 0 { 
							continue
						} else {
							break
						}
					} else {
						
						// 404 
						return nil, false, false 
					}
				
				}

			} 

			currentNode = nextNode
			continue

		}

	}

	handler, exists := currentNode.actions[method]

	if !exists {
		return nil, true, false
	} 

	return *handler, true, true

}


func(router *Router)ServeHTTP(w http.ResponseWriter, r *http.Request) {

	origin := r.Header.Get("Origin")

	// CORS
	w.Header().Set("Access-Control-Allow-Origin", router.cors.ValidateOrigin(origin))
	w.Header().Set("Access-Control-Allow-Methods", router.cors.AllowedMethods())
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

	if r.Method == "OPTIONS" {
		return 
	}

	// End CORS

	reqPath := r.URL.Path

	// global request log

	if origin == "" {
		origin = r.Host
	}

	log.Printf("%s:%s - %s", r.Method, reqPath, origin)

	rootHandler := &Root{
		writter: w,
		request: r,
		RequestParams: make(map[string]any),
	}

	// Search if the path exists
	handler, pathExists, methodExists := router.findHandler(reqPath, r.Method, rootHandler)	

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

	// Assign the global middlewares
	handler = router.SetGlobalMiddlewares(handler)
	handler(rootHandler)
}

// assign the list of middlewares to the handler
func(router *Router) SetGlobalMiddlewares(handler Handler) Handler {

	for _, midd := range *router.middlewares {
		handler = midd(handler)
	}

	return handler
}

// This function split the path and removes any empty value 
func(router *Router) explodePath(path string) []string {

	pathList := strings.Split(path, "/")

	var labels []string

	for _, str := range pathList {
		
		if str != "" {
			labels =append(labels, str)
		}

	}

	return labels
}

// Validate if an string is a path parameter like /user/{userId}

func(router *Router) isParameter(str string) bool {
	re := regexp.MustCompile("{([^}]+)}")
	match := re.MatchString(str)

	return match
}

// Remove the keys "{}" from a path parameter. clearParam("{userId}") returns "userId" 
func(router *Router) clearParam(str string) string {
	return strings.Trim(str, "{}")
}
