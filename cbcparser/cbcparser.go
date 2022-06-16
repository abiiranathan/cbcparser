package cbcparser

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrInvalidCSV       = errors.New("invalid csv file")
	ErrInsufficientRows = errors.New("csv must have at least 2 rows")
	ErrBlankCBCRecord   = errors.New("cbc record was a blank")
	ErrInvalidOutFormat = errors.New("invalid output format")
)

type OutFormat int

const (
	JSON OutFormat = iota
	JSONIndent
)

// Concrete type containing CBC fields values, units and flags
// Each CBCResult must implement a single Write method.
type CBCWriter interface {
	// Writes all the data into out writer.
	Write(out io.Writer, format OutFormat) error
}

// Data structure for multiple CBC records
type CBCMultiWriter []CBCWriter

type NormalRange struct {
	Lower float32 `json:"lower"`
	Upper float32 `json:"upper"`
}

type CBCValue struct {
	Value       float32     `json:"value"`
	Units       string      `json:"units"`
	Flag        string      `json:"flag"`
	NormalRange NormalRange `json:"normal_range"`
}

type CBCNormalRange struct {
	WBC        NormalRange `json:"wbc"`
	LYM        NormalRange `json:"lym"`
	MID        NormalRange `json:"mid"`
	GRA        NormalRange `json:"gra"`
	LYMPercent NormalRange `json:"lym_percent"`
	MIDPercent NormalRange `json:"mid_percent"`
	GRAPercent NormalRange `json:"gra_percent"`
	RBC        NormalRange `json:"rbc"`
	HGB        NormalRange `json:"hgb"`
	HCT        NormalRange `json:"hct"`
	MCV        NormalRange `json:"mcv"`
	MCH        NormalRange `json:"mch"`
	MCHC       NormalRange `json:"mchc"`
	RDWs       NormalRange `json:"rdw_s"`
	RDWc       NormalRange `json:"rdw_c"`
	PLT        NormalRange `json:"plt"`
	PCT        NormalRange `json:"pct"`
	MPV        NormalRange `json:"mpv"`

	// PDW is used by Edan machine
	PDW NormalRange `json:"pdw"`
	// Used by Human Machine
	PDWs NormalRange `json:"pdw_s"`
	// Used by Human Machine
	PDWc NormalRange `json:"pdw_c"`
	PLCC NormalRange `json:"plcc"`
	PLCR NormalRange `json:"plcr"`
}

// Parses a csv file with a single CBC record.
type CSVParser interface {
	Parse(r io.Reader, normal_ranges *CBCNormalRange) (CBCWriter, error)
}

// Parses a csv with multiple rows and returns a slice of CBCWriter structs.
type CSVMultiParser interface {
	ParseMulti(r io.Reader, normal_ranges *CBCNormalRange) (CBCMultiWriter, error)
}

func (list CBCMultiWriter) Write(out io.Writer, format OutFormat) error {
	var data []byte

	if format == JSON {
		d, err := json.Marshal(list)
		if err != nil {
			return err
		}
		data = d
	} else if format == JSONIndent {
		d, err := json.MarshalIndent(list, "", "	")
		if err != nil {
			return err
		}
		data = d
	} else {
		return ErrInvalidOutFormat
	}

	_, err := out.Write(data)
	return err
}

func ReadNormalRanges(r io.Reader) (*CBCNormalRange, error) {
	var normal_ranges CBCNormalRange
	err := json.NewDecoder(r).Decode(&normal_ranges)
	if err != nil {
		return nil, err
	}
	return &normal_ranges, nil
}
