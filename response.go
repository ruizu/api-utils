package apiutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, r interface{}, statusCode int) {
	b, _ := json.Marshal(r)
	if r.(Response).Callback != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "%s", string(b))
		return
	}
	w.Header().Set("Content-Type", "text/javascript")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s(%s)", r.(Response).Callback, string(b))
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int) {
	r := &Response{
		Errors: []ResponseError{
			{
				Code:   fmt.Sprintf("HTTP%d", statusCode),
				Title:  http.StatusText(statusCode),
				Detail: http.StatusText(statusCode),
			},
		},
	}
	b, _ := json.Marshal(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(b), statusCode)
}

type (
	Response struct {
		Callback string          `json:"-"`
		Links    ResponseLink    `json:"links,omitempty"`
		Meta     ResponseMeta    `json:"meta,omitempty"`
		Errors   []ResponseError `json:"errors,omitempty"`
	}
	ResponseMeta struct {
		ProcessTime float64 `json:"process_time"`
		TotalData   int64   `json:"total_data"`
	}
	ResponseLink struct {
		Self  string `json:"self,omitempty"`
		First string `json:"first,omitempty"`
		Last  string `json:"last,omitempty"`
		Next  string `json:"next,omitempty"`
		Prev  string `json:"prev,omitempty"`
	}
	ResponseError struct {
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}
)

func (r *Response) AddError(code, title, detail string) {
	e := ResponseError{
		Code: code,
		Title: title,
		Detail: detail,
	}
	r.Errors = append(r.Errors, e)
}

func (r *Response) ClearError() {
	r.Errors = make([]ResponseError, 0)
}
