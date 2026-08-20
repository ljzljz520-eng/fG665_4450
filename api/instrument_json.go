package api

import (
	"encoding/json"
	"instrumentarchive/model"
	"net/http"
)

type InstrumentPayload struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Laboratory   string `json:"laboratory"`
	Owner        string `json:"owner"`
	PurchaseDate string `json:"purchase_date"`
}

func (p InstrumentPayload) Instrument(id string) model.Instrument {
	return model.NewInstrument(id, p.Number, p.Name, p.Laboratory, p.Owner, p.PurchaseDate)
}
func DecodePayload(r *http.Request) (InstrumentPayload, error) {
	var p InstrumentPayload
	e := json.NewDecoder(r.Body).Decode(&p)
	return p, e
}
func EncodeInstrument(w http.ResponseWriter, i model.Instrument) { writeJSON(w, 200, i) }
func EncodeDetail(w http.ResponseWriter, d model.Detail)         { writeJSON(w, 200, d) }
func EncodeList(w http.ResponseWriter, items []model.Instrument) { writeJSON(w, 200, items) }
func RequestID(r *http.Request) string                           { return r.Header.Get("X-Request-ID") }
