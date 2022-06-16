package edan

import "github.com/abiiranathan/cbcparser/cbcparser"

// Sample ID,Mode,Analysis Time,WBC(10^3/��L),LYM#(10^3/��L),LYM%(%),MXD#(),MXD%(),NEUT#(),NEUT%(),RBC(10^6/��L),HGB(g/dL),HCT(%),MCV(fL),MCH(pg),MCHC(g/dL),RDW_CV(%),RDW_SD(fL),PLT(10^3/��L),PDW(fL),MPV(fL),PCT(%),P_LCR(%),P_LCC(10^3/��L)
// Structure to store data parsed from the text file
// exported by the Edan Pro30 CBC Machine.
type edanCBCResult struct {
	SID          string `json:"sid"`
	Mode         string `json:"mode"`
	AnalysisTime string `json:"analysis_time"`
	PID          string `json:"pid"`

	WBC        cbcparser.CBCValue `json:"wbc"`
	LYM        cbcparser.CBCValue `json:"lym"`
	LYMPercent cbcparser.CBCValue `json:"lym_percent"`
	MID        cbcparser.CBCValue `json:"mid"`
	MIDPercent cbcparser.CBCValue `json:"mid_percent"`
	GRA        cbcparser.CBCValue `json:"gra"`
	GRAPercent cbcparser.CBCValue `json:"gra_percent"`

	RBC  cbcparser.CBCValue `json:"rbc"`
	HGB  cbcparser.CBCValue `json:"hgb"`
	HCT  cbcparser.CBCValue `json:"hct"`
	MCV  cbcparser.CBCValue `json:"mcv"`
	MCH  cbcparser.CBCValue `json:"mch"`
	MCHC cbcparser.CBCValue `json:"mchc"`
	RDWc cbcparser.CBCValue `json:"rdw_c"`
	RDWs cbcparser.CBCValue `json:"rdw_s"`

	PLT cbcparser.CBCValue `json:"plt"`
	PDW cbcparser.CBCValue `json:"pdw"`
	MPV cbcparser.CBCValue `json:"mpv"`
	PCT cbcparser.CBCValue `json:"pct"`

	PLCC cbcparser.CBCValue `json:"plcc"`
	PLCR cbcparser.CBCValue `json:"plcr"`
}

type edanCBCResultMulti []edanCBCResult
