package api

import (
	"instrumentarchive/model"
	"net/http"
	"strings"
)

func ValidateMethod(r *http.Request, allowed ...string) bool {
	for _, m := range allowed {
		if r.Method == m {
			return true
		}
	}
	return false
}
func InstrumentFromQuery(r *http.Request) model.Instrument {
	q := r.URL.Query()
	return model.NewInstrument(q.Get("id"), q.Get("number"), q.Get("name"), q.Get("laboratory"), q.Get("owner"), q.Get("purchase_date"))
}
func HeaderRole(r *http.Request) string {
	role := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Role")))
	if role == "admin" {
		return role
	}
	return "viewer"
}
func IsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
func SetCache(w http.ResponseWriter, seconds int) {
	w.Header().Set("Cache-Control", "max-age="+itoa(seconds))
}
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	digits := ""
	for v > 0 {
		digits = string(rune('0'+v%10)) + digits
		v /= 10
	}
	return digits
}
