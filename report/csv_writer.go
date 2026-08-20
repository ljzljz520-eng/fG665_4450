package report

import (
	"encoding/csv"
	"instrumentarchive/model"
	"io"
	"strconv"
)

func WriteCSV(w io.Writer, items []model.Instrument) error {
	cw := csv.NewWriter(w)
	if e := cw.Write([]string{"id", "number", "name", "laboratory", "owner", "purchase_date", "status"}); e != nil {
		return e
	}
	for _, i := range items {
		if e := cw.Write([]string{i.ID, i.Number, i.Name, i.Laboratory, i.Owner, i.PurchaseDate, string(i.Status)}); e != nil {
			return e
		}
	}
	cw.Flush()
	return cw.Error()
}
func EncodeRows(items []model.Instrument) [][]string {
	rows := [][]string{}
	for _, i := range items {
		rows = append(rows, []string{i.ID, i.Number, i.Name, i.Laboratory, i.Owner, i.PurchaseDate, string(i.Status), strconv.FormatBool(i.Archived)})
	}
	return rows
}
func Headers() []string {
	return []string{"id", "number", "name", "laboratory", "owner", "purchase_date", "status", "archived"}
}
func ParseStatusColumn(raw string) model.Status { return model.NormalizeStatus(raw) }
