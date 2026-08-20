package report

import (
	"encoding/csv"
	"fmt"
	"instrumentarchive/model"
	"io"
	"strings"
)

type Row struct {
	Instrument model.Instrument
	Error      string
}
type Summary struct {
	Total    int
	Imported int
	Rejected int
	Errors   []string
}

func ParseCSV(r io.Reader) ([]Row, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true
	header, e := cr.Read()
	if e != nil {
		return nil, e
	}
	if len(header) < 5 {
		return nil, fmt.Errorf("header requires 5 columns")
	}
	rows := []Row{}
	for {
		rec, e := cr.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		if len(rec) < 5 {
			rows = append(rows, Row{Error: "columns"})
			continue
		}
		i := model.NewInstrument(strings.TrimSpace(rec[0]), strings.TrimSpace(rec[1]), strings.TrimSpace(rec[2]), strings.TrimSpace(rec[3]), strings.TrimSpace(rec[4]), "")
		if e = i.Validate(); e != nil {
			rows = append(rows, Row{Instrument: i, Error: e.Error()})
		} else {
			rows = append(rows, Row{Instrument: i})
		}
	}
	return rows, nil
}
func Summarize(rows []Row) Summary {
	s := Summary{Total: len(rows), Errors: []string{}}
	for _, r := range rows {
		if r.Error != "" {
			s.Rejected++
			s.Errors = append(s.Errors, r.Error)
		} else {
			s.Imported++
		}
	}
	return s
}
func Render(s Summary) string {
	return fmt.Sprintf("total=%d imported=%d rejected=%d", s.Total, s.Imported, s.Rejected)
}
func ValidRows(rows []Row) []model.Instrument {
	out := []model.Instrument{}
	for _, r := range rows {
		if r.Error == "" {
			out = append(out, r.Instrument)
		}
	}
	return out
}
