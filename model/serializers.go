package model

import (
	"encoding/json"
	"errors"
	"strings"
)

func EncodeInstrument(i Instrument) ([]byte, error) { return json.Marshal(i) }
func DecodeInstrument(data []byte) (Instrument, error) {
	var i Instrument
	if len(data) == 0 {
		return i, errors.New("empty instrument")
	}
	e := json.Unmarshal(data, &i)
	return NormalizeInstrument(i), e
}
func EncodeCalibration(c Calibration) ([]byte, error) { return json.Marshal(c) }
func DecodeCalibration(data []byte) (Calibration, error) {
	var c Calibration
	e := json.Unmarshal(data, &c)
	return c, e
}
func ParseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "y", "是":
		return true
	}
	return false
}
func FormatBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
func JoinLabels(labels []string) string { return strings.Join(labels, "、") }
func SplitLabels(raw string) []string {
	parts := strings.Split(raw, "、")
	out := []string{}
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, strings.TrimSpace(p))
		}
	}
	return out
}
