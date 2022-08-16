package api_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/TDD-all-the-things/learn-by-testing/api"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestHello(t *testing.T) {

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
