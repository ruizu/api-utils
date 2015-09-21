package utils

import (
	"log"
	"net/http"
)

var NotFoundHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	WriteErrorResponse(w, http.StatusNotFound)
}

var MethodNotAllowedHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed)
}

var PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Print("PANIC: ", err)
	WriteErrorResponse(w, http.StatusInternalServerError)
}
