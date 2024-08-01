package helpers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

const (
	UserEmail    = "test@test.com"
	UserPassword = "00000000"
)

// Test defines a structure for specifying input and output data of a single test case.
type Test struct {
	Description string

	// Test input
	Route   string
	Method  string
	Body    io.Reader
	Headers []Header

	// Check
	CheckError bool
	CheckBody  bool
	CheckCode  bool

	// Expected output
	ExpectedError bool
	ExpectedCode  int
	ExpectedBody  string
}

// Header represents an header value.
type Header struct {
	Key   string
	Value string
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)

	return rr
}

// JsonToString converts a JSON to a string.
func JsonToString(d interface{}) string {
	b, err := json.Marshal(d)
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return string(b)
}
