package auth

import "errors"

type Role string

const (
	Admin  Role = "admin"
	Viewer Role = "viewer"
)

type User struct {
	ID   string
	Name string
	Role Role
}

func NewUser(id, name string, r Role) User { return User{ID: id, Name: name, Role: r} }
func Authorize(u User, action string) error {
	if u.ID == "" {
		return errors.New("user required")
	}
	if action == "view" {
		return nil
	}
	if u.Role != Admin {
		return errors.New("admin required")
	}
	return nil
}
func CanWrite(u User) bool { return u.Role == Admin }
func CanRead(u User) bool  { return u.Role == Admin || u.Role == Viewer }
func ParseRole(raw string) Role {
	if raw == "admin" {
		return Admin
	}
	return Viewer
}
