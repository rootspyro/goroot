package goroot

import (
	"encoding/json"
	"log"
	"net/http"
)

// Root it's the kernel of the handlers of GoRoot

type Root struct {
	Writter http.ResponseWriter
	Request *http.Request
}

type Handler func( root *Root )

func(root *Root)parseJson(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func(root *Root)serverERROR(err error) {
	log.Printf("%v", err)
	root.Writter.WriteHeader(http.StatusInternalServerError)
}

func(root *Root)Json(data interface {}) {

	payload, err := root.parseJson(data)

	if err != nil {
		root.serverERROR(err)
		return
	}

	root.Writter.Header().Set("Content-Type", "application/json")
	root.Writter.Write(payload)
}

