package goroot

import (
	"encoding/json"
	"log"
	"net/http"
)

// Root it's the kernel of the handlers of GoRoot

type Root struct {
	writter http.ResponseWriter
	request *http.Request
	_status int
}

type Handler func( root *Root )

// STATUS

// this function allows the programmer to easily write an http status
func(root *Root)Status(code int) *Root {

	root._status = code
	return root
}

// Default Statuses

// 2XX List
func(root *Root)OK() *Root {
	root.Status(200)
	return root
}

func(root *Root)Created() *Root {
	root.Status(201)
	return root
}

// 4XX list

func(root *Root)BadRequest() *Root {
	root.Status(400)
	return root
}

func(root *Root)Unauthorized() *Root {
	root.Status(401)
	return root
}

func(root *Root)Forbidden() *Root {
	root.Status(403)
	return root
}

func(root *Root)NotFound() *Root {
	root.Status(404)
	return root
}

func(root *Root)MethodNotAllowed() *Root {
	root.Status(405)
	return root
}

// 5XX List

func(root *Root)InternalServerError() *Root {
	root.Status(500)
	return root
}

func(root *Root)NotImplemented() *Root {
	root.Status(501)
	return root
}


// RESPONSES

func(root *Root)Send(data string) {

	httpCode := 200 

	if root._status > 0 {
		
		httpCode = root._status

	}

	root.writter.WriteHeader(httpCode)
	root.writter.Write([]byte(data))
}

func(root *Root)parseJson(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func(root *Root)Json(data interface {}) {

	payload, err := root.parseJson(data)

	if err != nil {
		log.Printf("%v", err)
		root.Status(500)
	}

	httpCode := 200 

	if root._status > 0 {
		
		httpCode = root._status

	}

	root.writter.Header().Set("Content-Type", "application/json")
	root.writter.WriteHeader(httpCode)
	root.writter.Write(payload)
}

