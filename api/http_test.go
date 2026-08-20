package api

import (
	"instrumentarchive/service"
	"instrumentarchive/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHTTPList(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	rr := httptest.NewRecorder()
	New(service.New(s, service.StaticClock{"t"})).Routes().ServeHTTP(rr, httptest.NewRequest("GET", "/health", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
}
