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


// RESPONSES

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

