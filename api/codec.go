package api

import (
	"encoding/json"
	"instrumentarchive/model"
	"net/http"
)

func DecodeInstrument(r *http.Request) (model.Instrument, error) {
	var i model.Instrument
	err := json.NewDecoder(r.Body).Decode(&i)
	return i, err
}
func QueryFilter(r *http.Request) model.Filter {
	q := r.URL.Query()
	return model.Filter{Query: q.Get("q"), Laboratory: q.Get("laboratory"), Status: model.NormalizeStatus(q.Get("status"))}
}
func SetRole(r *http.Request) string {
	if r.Header.Get("X-Role") == "admin" {
		return "admin"
	}
	return "viewer"
}
