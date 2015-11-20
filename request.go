package apiutils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"net/http"
	"net/url"
)

type Request struct {
	UserID     int64
	Sort       []string
	Filter     map[string][]string
	Device     string
	Callback   string
	PageNumber int
	PageSize   int
	PageLimit  int
	PageOffset int
}

var (
	validVariableName *regexp.Regexp
	validFilterName   *regexp.Regexp
)

func init() {
	validVariableName = regexp.MustCompile("^[a-zA-Z\\$_]+[a-zA-Z0-9\\$_]*(\\.[a-zA-Z\\$_]+[a-zA-Z0-9\\$_]*)*$")
	validFilterName = regexp.MustCompile("^filter\\[([^\\]]+?)\\]$")
}

func ParseRequest(r *http.Request) (Request, error) {
	req := Request{}

	req.UserID, _ = strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
	req.Device = r.FormValue("device")
	req.Callback = r.FormValue("callback")
	if req.Callback != "" && !validVariableName.Match([]byte(req.Callback)) {
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

	req.PageLimit = req.PageSize
	req.PageOffset = req.PageSize*(req.PageNumber-1)
	req.Filter = parseRequestFilters(r.Form)

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

func parseRequestFilters(form url.Values) map[string][]string {
	filters := make(map[string][]string)
	for k, v := range form {
		if v[0] == "" {
			continue
		}
		if t := validFilterName.FindStringSubmatch(k); len(t) > 0 {
			filters[t[1]] = strings.Split(v[0], ",")
		}
	}
	return filters
}

func (req Request) FilterString(name string) string {
	if v, ok := req.Filter[name]; ok {
		return v[0]
	}
	return ""
}

func (req Request) FilterInt(name string) int {
	v := req.FilterString(name)
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}

func (req Request) FilterBool(name string) bool {
	v := req.FilterInt(name)
	if v > 0 {
		return true
	}
	return false
}

func ParseBodyRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

func GetIDs(value string) ([]int64, error) {
	idn := []int64{}
	ida := strings.Split(value, ",")
	if len(ida) == 1 && ida[0] == "" {
		return idn, nil
	}

	seen := map[string]bool{}
	for _, v := range ida {
		if _, ok := seen[v]; ok {
			continue
		}

		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		idn = append(idn, n)
		seen[v] = true
	}

	return idn, nil
}

// deprecated
func GetID(value string) ([]int64, error) {
	return GetIDs(value)
}
