package apiutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"net/http"
)

type Request struct {
	Sort       []string
	Device     string
	Callback   string
	PageNumber int
	PageSize   int
}

var MaximumPageSize int = 25
var validVariableName = regexp.Compile("^[a-zA-Z\\$_]+[a-zA-Z0-9\\$_]*(\\.[a-zA-Z\\$_]+[a-zA-Z0-9\\$_]*)*$")

func ParseRequest(r *http.Request) (Request, error) {
	req := Request{}

	req.Device = r.FormValue("device")
	req.Callback = r.FormValue("callback")
	if req.Callback && !validVariableName.Match([]byte(req.Callback)) {
		return Request{}, fmt.Errorf("invalid callback")
	}

	req.PageNumber, _ = strconv.Atoi(r.FormValue("page[number]"))
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}

	req.PageSize, _ = strconv.Atoi(r.FormValue("page[size]"))
	if req.PageSize < 1 {
		req.PageSize = 1
	}

	if req.PageSize > MaximumPageSize {
		return Request{}, fmt.Errorf("invalid page[size]")
	}

	req.Sort = []string{}
	sorts := strings.Split(r.FormValue("sort"), ",")
	if len(sorts) == 1 && sorts[0] == "" {
		return req, nil
	}

	seen := map[string]bool{}
	for _, v := range sorts {
		if _, ok := seen[v]; ok {
			continue
		}
		req.Sort = append(req.Sort, v)
		seen[v] = true
	}

	return req, nil
}

func GetID(value string) ([]int, error) {
	idn := []int{}
	ida := strings.Split(value, ",")
	if len(ida) == 1 && ida[0] == "" {
		return idn, nil
	}

	seen := map[string]bool{}
	for _, v := range ida {
		if _, ok := seen[v]; ok {
			continue
		}

		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		idn = append(idn, n)
		seen[v] = true
	}

	return idn, nil
}
