package app

import (
	"context"
	"net/http"

	"rendering/internal/domain"
)

// Pages selects the template name and view data for a request.
type Pages struct{}

// NewPages returns the page use case.
func NewPages() *Pages {
	return &Pages{}
}

// Index returns the default home page.
func (*Pages) Index(context.Context) (domain.Page, error) {
	return domain.Page{
		Name:   "index",
		Status: http.StatusOK,
		View:   domain.IndexView(),
	}, nil
}

// Status maps a code to an error template. Unknown codes fall back to 500.
func (*Pages) Status(_ context.Context, code int) (domain.Page, error) {
	name, ok := statusNames[code]
	if !ok {
		code = http.StatusInternalServerError
		name = statusNames[code]
	}
	return domain.Page{
		Name:   name,
		Status: code,
		View:   domain.StatusView(code),
	}, nil
}

var statusNames = map[int]string{
	http.StatusBadRequest:          "400",
	http.StatusUnauthorized:        "401",
	http.StatusForbidden:           "403",
	http.StatusNotFound:            "404",
	http.StatusMethodNotAllowed:    "405",
	http.StatusInternalServerError: "500",
	http.StatusBadGateway:          "502",
	http.StatusServiceUnavailable:  "503",
}
