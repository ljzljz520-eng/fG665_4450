package api

import (
	"encoding/json"
	"instrumentarchive/model"
	"net/http"
	"strconv"
)

type Page struct {
	Items []model.Instrument `json:"items"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
	Total int                `json:"total"`
}

func Paginate(items []model.Instrument, page, size int) Page {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start > len(items) {
		start = len(items)
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return Page{Items: items[start:end], Page: page, Size: size, Total: len(items)}
}
func DecodePage(r *http.Request) (int, int) {
	q := r.URL.Query()
	p, _ := strconv.Atoi(q.Get("page"))
	s, _ := strconv.Atoi(q.Get("size"))
	return p, s
}
func Encode(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
func ErrorPayload(e error) map[string]string {
	if e == nil {
		return map[string]string{}
	}
	return map[string]string{"error": e.Error()}
}
func StatusCode(err error) int {
	if err == nil {
		return 200
	}
	return 400
}
func ApplyPage(items []model.Instrument, r *http.Request) Page {
	p, s := DecodePage(r)
	return Paginate(items, p, s)
}
