package auth

import "testing"

func TestAuthorize(t *testing.T) {
	u := NewUser("1", "A", Admin)
	if Authorize(u, "archive") != nil {
		t.Fatal("admin")
	}
	v := NewUser("2", "V", Viewer)
	if Authorize(v, "archive") == nil || Authorize(v, "view") != nil {
		t.Fatal("viewer")
	}
}
