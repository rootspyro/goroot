package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseHandler struct{
	w http.ResponseWriter
}

func NewResponseHandler(w http.ResponseWriter) *ResponseHandler {
	return &ResponseHandler{
		w: w,
	}
}

func(rh *ResponseHandler)ParseJson(data interface{}) ([]byte, error) {

	return json.Marshal(data)

}

func(rh *ResponseHandler)ServerERROR(){
	rh.w.WriteHeader(http.StatusInternalServerError)
}

func(rh *ResponseHandler)JsonNotFound( data interface{} ) {

	payload, err := rh.ParseJson(data)
	
	if err != nil {

		log.Printf("%v", err)
		rh.ServerERROR()
		return

	}

	rh.w.Header().Set("Content-Type","application/json")
	rh.w.WriteHeader(http.StatusNotFound)
	rh.w.Write(payload)
}

func(rh *ResponseHandler)JsonInternalERROR( data interface{} ) {

	payload, err := rh.ParseJson(data)
	
	if err != nil {

		log.Printf("%v", err)
		rh.ServerERROR()
		return

	}

	rh.w.Header().Set("Content-Type","application/json")
	rh.w.WriteHeader(http.StatusInternalServerError)
	rh.w.Write(payload)
}


func(rh *ResponseHandler)JsonOK( data interface{}) {

	payload, err := rh.ParseJson(data)
	
	if err != nil {

		log.Printf("%v", err)
		rh.ServerERROR()
		return

	}

	rh.w.Header().Set("Content-Type","application/json")
	rh.w.WriteHeader(http.StatusOK)
	rh.w.Write(payload)
}
