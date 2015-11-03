package apiutils

import (
	"log"
	"net/http"
)

var Debug bool = false

var NotFoundHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	WriteErrorResponse(w, http.StatusNotFound)
}

var MethodNotAllowedHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed)
}

var PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
	if Debug {
		r := &Response{}
		r.AddError(
			fmt.Sprintf("HTTP%d", http.StatusInternalServerError),
			http.StatusInternalServerError,
			err)
		WriteResponse(w, r, http.StatusInternalServerError)
		return
	}
	WriteErrorResponse(w, http.StatusInternalServerError)
}
