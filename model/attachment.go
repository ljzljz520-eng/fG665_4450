package model

import (
	"errors"
	"path/filepath"
	"strings"
)

func (a Attachment) Validate() error {
	if a.ID == "" {
		return errors.New("attachment id required")
	}
	if a.InstrumentID == "" {
		return errors.New("instrument id required")
	}
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("attachment name required")
	}
	if filepath.Base(a.Name) != a.Name {
		return errors.New("attachment name must be basename")
	}
	if len(a.Data) == 0 {
		return errors.New("attachment data required")
	}
	return nil
}
func (a Attachment) Extension() string { return strings.ToLower(filepath.Ext(a.Name)) }
func (a Attachment) IsDocument() bool {
	switch a.Extension() {
	case ".pdf", ".doc", ".docx", ".txt":
		return true
	}
	return false
}
func (a Attachment) IsImage() bool {
	switch a.Extension() {
	case ".png", ".jpg", ".jpeg", ".gif":
		return true
	}
	return false
}
func (a Attachment) Size() int       { return len(a.Data) }
func (a Attachment) Summary() string { return a.Name + " (" + a.ContentType + ")" }
