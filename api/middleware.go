package api

import (
	"net/http"
	"time"
)

type AccessLog struct {
	Method   string
	Path     string
	Status   int
	Duration time.Duration
}
type Recorder struct{ Logs []AccessLog }

func NewRecorder() *Recorder { return &Recorder{Logs: []AccessLog{}} }
func (r *Recorder) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, req)
		r.Logs = append(r.Logs, AccessLog{Method: req.Method, Path: req.URL.Path, Status: rw.status, Duration: time.Since(start)})
	})
}
func (r *Recorder) Count(status int) int {
	n := 0
	for _, l := range r.Logs {
		if l.Status == status {
			n++
		}
	}
	return n
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(s int)           { w.status = s; w.ResponseWriter.WriteHeader(s) }
func (w *statusWriter) Write(b []byte) (int, error) { return w.ResponseWriter.Write(b) }
func WithJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Archive-Version", "1")
		next.ServeHTTP(w, r)
	})
}
