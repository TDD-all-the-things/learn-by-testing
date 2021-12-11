package api_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gojustforfun/learn-by-test/api"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
)

// TestMySuite is the go test entry point.
func TestAPISuite(t *testing.T) {
	tdsuite.Run(t, &APISuite{})
}

type APISuite struct {
	ta  *tdhttp.TestAPI
	mux *http.ServeMux
}

func (s *APISuite) Setup(t *td.T) error {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/hello", api.Hello)
	fmt.Println("Setup ")
	return nil
}

// Destroy is called after all tests are run.
// Destroy is not called if Setup returned an error.
func (s *APISuite) Destroy(t *td.T) error {
	s.mux = nil
	s.ta = nil
	fmt.Println("Destroy ")
	return nil
}

func (s *APISuite) PreTest(t *td.T, testName string) error {
	fmt.Println("PreTest")
	return nil
}

func (s *APISuite) PostTest(t *td.T, testName string) error {
	fmt.Println("PostTest")
	return nil
}

func (s *APISuite) TestHello(t *td.T) {

	ta := tdhttp.NewTestAPI(t, http.HandlerFunc(api.Hello))

	ta.Run("/GET", func(t *tdhttp.TestAPI) {
		t.Get("/hello").
			CmpStatus(http.StatusMethodNotAllowed).
			CmpHeader(http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}).
			CmpJSONBody(td.JSON(`{"errno":1, "msg":"Method not allowed", "data":{}}`))
	})

	ta.Run("/POST Form", func(t *tdhttp.TestAPI) {
		t.PostForm("/hello", url.Values{"name": []string{"Longyue"}}).
			CmpStatus(http.StatusOK).
			CmpHeader(http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}).
			CmpJSONBody(td.JSON(`{"errno":0, "msg":"Hello Longyue", "data":{}}`))
	})

	ta.Run("/POST", func(t *tdhttp.TestAPI) {
		t.Post("/hello", strings.NewReader(`name=Longyue`), "Content-Type", "application/x-www-form-urlencoded").
			CmpHeader(http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}).
			CmpStatus(http.StatusOK).
			CmpJSONBody(td.JSON(`{"errno":0, "msg":"Hello Longyue", "data":{}}`))
	})
}
