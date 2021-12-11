package api

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"errno":1, "msg":"Method not allowed", "data":{}}`))
		return
	}

	r.ParseForm()
	name := r.Form.Get("name")
	w.Write([]byte(fmt.Sprintf(`{"errno":0, "msg":"Hello %s", "data":{}}`, name)))
}
