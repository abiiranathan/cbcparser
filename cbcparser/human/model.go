package human

import "github.com/abiiranathan/cbcparser/cbcparser"

// Structure to store data parsed from the text file
// exported by the HUMAN 30 CBC Machine.
type humanCBCResult struct {
	SampleID  string `json:"sample_id"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	PatientID string `json:"patient_id"`
	BirthDate string `json:"birth_date"`

	WBC        cbcparser.CBCValue `json:"wbc"`
	LYM        cbcparser.CBCValue `json:"lym"`
	MID        cbcparser.CBCValue `json:"mid"`
	GRA        cbcparser.CBCValue `json:"gra"`
	LYMPercent cbcparser.CBCValue `json:"lym_percent"`
	MIDPercent cbcparser.CBCValue `json:"mid_percent"`
	GRAPercent cbcparser.CBCValue `json:"gra_percent"`
	RBC        cbcparser.CBCValue `json:"rbc"`
	HGB        cbcparser.CBCValue `json:"hgb"`
	HCT        cbcparser.CBCValue `json:"hct"`
	MCV        cbcparser.CBCValue `json:"mcv"`
	MCH        cbcparser.CBCValue `json:"mch"`
	MCHC       cbcparser.CBCValue `json:"mchc"`
	RDWs       cbcparser.CBCValue `json:"rdw_s"`
	RDWc       cbcparser.CBCValue `json:"rdw_c"`
	PLT        cbcparser.CBCValue `json:"plt"`
	PCT        cbcparser.CBCValue `json:"pct"`
	MPV        cbcparser.CBCValue `json:"mpv"`
	PDWs       cbcparser.CBCValue `json:"pdw_s"`
	PDWc       cbcparser.CBCValue `json:"pdw_c"`
	PLCC       cbcparser.CBCValue `json:"plcc"`
	PLCR       cbcparser.CBCValue `json:"plcr"`

	Type    string `json:"type"`
	Warning string `json:"warning"`
}

type humanCBCResultMulti []humanCBCResult
