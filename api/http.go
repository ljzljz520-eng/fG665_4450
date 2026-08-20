package api

import (
	"encoding/json"
	"instrumentarchive/auth"
	"instrumentarchive/model"
	"instrumentarchive/service"
	"net/http"
	"strings"
)

type Server struct{ Service *service.Service }

func New(s *service.Service) *Server { return &Server{Service: s} }
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/instruments", s.instruments)
	mux.HandleFunc("/instruments/", s.instrument)
	return mux
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
func (s *Server) instruments(w http.ResponseWriter, r *http.Request) {
	role := auth.ParseRole(r.Header.Get("X-Role"))
	if r.Method == http.MethodGet {
		q := r.URL.Query()
		items, e := s.Service.Search(model.Filter{Query: q.Get("q"), Laboratory: q.Get("laboratory"), Status: model.NormalizeStatus(q.Get("status")), IncludeArchived: role == auth.Admin}, string(role))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, items)
		return
	}
	if r.Method == http.MethodPost {
		if role != auth.Admin {
			writeErrStatus(w, http.StatusForbidden, "admin required")
			return
		}
		var i model.Instrument
		if e := json.NewDecoder(r.Body).Decode(&i); e != nil {
			writeErr(w, e)
			return
		}
		if e := s.Service.Create(i, string(role)); e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, http.StatusCreated, i)
		return
	}
	writeErrStatus(w, 405, "method not allowed")
}
func (s *Server) instrument(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/instruments/")
	role := auth.ParseRole(r.Header.Get("X-Role"))
	if r.Method == http.MethodGet {
		d, e := s.Service.Detail(id, string(role))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, d)
		return
	}
	writeErrStatus(w, 405, "method not allowed")
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, e error) { writeErrStatus(w, 400, e.Error()) }
func writeErrStatus(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
