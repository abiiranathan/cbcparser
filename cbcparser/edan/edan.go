// parses cbc tsv/csv to print the cbc report to stdout
// Ascii char go from:
// https://theasciicode.com.ar/extended-ascii-code/box-drawing-character-single-line-upper-left-corner-ascii-code-218.html
//
// CBC format:
// Sample ID,Mode,Analysis Time,WBC(10^3/��L),LYM#(10^3/��L),LYM%(%),MXD#(),MXD%(),NEUT#(),NEUT%(),RBC(10^6/��L),HGB(g/dL),HCT(%),MCV(fL),MCH(pg),MCHC(g/dL),RDW_CV(%),RDW_SD(fL),PLT(10^3/��L),PDW(fL),MPV(fL),PCT(%),P_LCR(%),P_LCC(10^3/��L)
package edan

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/abiiranathan/cbcparser/cbcparser"
)

var (
	// Regular expression to extract units from brackets
	// e.g WBC(10^3/uL) extracts 10^3/uL
	UnitsRegex = regexp.MustCompile(`\((.*?)\)`)
)

const (
	nfields   = 24
	separator = ','
)

func parse_float(value string) float32 {
	fvalue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0.0
	}
	return float32(fvalue)
}

// Extracts units from header, replacing invalid unicode
// with μ(micro) symbol.
func extract_units(value string) string {
	rs := UnitsRegex.FindStringSubmatch(strings.ToValidUTF8(value, "μ"))

	if len(rs) == 2 {
		return rs[1]
	}

	return ""
}

// Returns L or H if value is out of range or an empty string.
func get_flag(value float32, nrange cbcparser.NormalRange) string {
	if value < nrange.Lower {
		return "L"
	}

	if value > nrange.Upper {
		return "H"
	}

	return ""
}

// Initialize a new CSV Parser
func New() cbcparser.CSVParser {
	return &EdanCBCResult{}
}

func NewMultiParser() cbcparser.CSVMultiParser {
	return &EdanCBCResultMulti{}
}

func (cbc EdanCBCResult) Write(out io.Writer, format cbcparser.OutFormat) error {
	var data []byte

	if format == cbcparser.JSON {
		d, err := json.Marshal(cbc)
		if err != nil {
			return err
		}
		data = d
	} else if format == cbcparser.JSONIndent {
		d, err := json.MarshalIndent(cbc, "", "   ")
		if err != nil {
			return err
		}
		data = d
	} else {
		return cbcparser.ErrInvalidOutFormat
	}

	_, err := out.Write(data)
	return err
}

// ParseTSV reads from r and parses the data into an slice of CBCResult structs.
// The text file is expected to be in the csv format.
// The first line of the file is expected to be the header.
//
// The header is expected to be in the following format(with 24 columns):
func (cbc EdanCBCResult) Parse(r io.Reader, normal_ranges *cbcparser.CBCNormalRange) (cbcparser.CBCWriter, error) {
	reader := csv.NewReader(r)
	reader.Comma = separator // ',' or '\t'
	reader.FieldsPerRecord = nfields

	data, err := reader.ReadAll()

	if err != nil {
		return nil, cbcparser.ErrInvalidCSV
	}

	if len(data) < 2 {
		return nil, cbcparser.ErrInsufficientRows
	}

	// initialize with numrows excluding header row
	headers := data[0]
	row := data[1]
	result := set_cbc_value(headers, row, normal_ranges)
	return result, nil
}

// MultiParse reads from r and parses the data into an slice of a CBCWriter struct.
func (EdanCBCResultMulti) ParseMulti(r io.Reader, normal_ranges *cbcparser.CBCNormalRange) (cbcparser.CBCMultiWriter, error) {
	reader := csv.NewReader(r)
	reader.Comma = separator // ',' or '\t'
	reader.FieldsPerRecord = nfields

	data, err := reader.ReadAll()

	if err != nil {
		return nil, cbcparser.ErrInvalidCSV
	}

	if len(data) < 2 {
		return nil, cbcparser.ErrInsufficientRows
	}

	headers := data[0]
	var result []cbcparser.CBCWriter
	for i := 1; i < len(data); i++ {
		row := data[i]
		result = append(result, set_cbc_value(headers, row, normal_ranges))
	}

	return result, nil
}

/*
row index are as below.
-------------------------
Indexes: 0 - 4 are identifiers
Indexes: 5 - 48 have units and flags
Index: 49 (Type)
Index: 50 (Warning)
*/
func set_cbc_value(headers []string, row []string, normal_ranges *cbcparser.CBCNormalRange) cbcparser.CBCWriter {
	cbcRes := EdanCBCResult{}
	// Patient Identifiers
	cbcRes.SID = row[0]
	cbcRes.Mode = row[1]
	cbcRes.AnalysisTime = row[2]

	// Absoluet counts
	wbc := parse_float(row[3])
	cbcRes.WBC = cbcparser.CBCValue{
		Value: wbc,
		Units: extract_units(headers[3]),
	}

	lym := parse_float(row[4])
	cbcRes.LYM = cbcparser.CBCValue{
		Value: lym,
		Units: extract_units(headers[4]),
	}

	lym_percent := parse_float(row[5])
	cbcRes.LYMPercent = cbcparser.CBCValue{
		Value: lym_percent,
		Units: extract_units(headers[5]),
	}

	mid := parse_float(row[6])
	cbcRes.MID = cbcparser.CBCValue{
		Value: mid,
		Units: extract_units(headers[6]),
	}

	mid_percent := parse_float(row[7])
	cbcRes.MIDPercent = cbcparser.CBCValue{
		Value: mid_percent,
		Units: extract_units(headers[7]),
	}

	granulocytes := parse_float(row[8])
	cbcRes.GRA = cbcparser.CBCValue{
		Value: granulocytes,
		Units: extract_units(headers[8]),
	}

	gra_percent := parse_float(row[9])
	cbcRes.GRAPercent = cbcparser.CBCValue{
		Value: gra_percent,
		Units: extract_units(headers[9]),
	}

	rbc := parse_float(row[10])
	cbcRes.RBC = cbcparser.CBCValue{
		Value: rbc,
		Units: extract_units(headers[10]),
	}

	hgb := parse_float(row[11])
	cbcRes.HGB = cbcparser.CBCValue{
		Value: hgb,
		Units: extract_units(headers[11]),
	}

	hct := parse_float(row[12])
	cbcRes.HCT = cbcparser.CBCValue{
		Value: hct,
		Units: extract_units(headers[12]),
	}

	mcv := parse_float(row[13])
	cbcRes.MCV = cbcparser.CBCValue{
		Value: mcv,
		Units: extract_units(headers[13]),
	}

	mch := parse_float(row[14])
	cbcRes.MCH = cbcparser.CBCValue{
		Value: mch,
		Units: extract_units(headers[14]),
	}

	mchc := parse_float(row[15])
	cbcRes.MCHC = cbcparser.CBCValue{
		Value: mchc,
		Units: extract_units(headers[15]),
	}

	rdwc := parse_float(row[16])
	cbcRes.RDWc = cbcparser.CBCValue{
		Value: rdwc,
		Units: extract_units(headers[16]),
	}

	rdws := parse_float(row[17])
	cbcRes.RDWs = cbcparser.CBCValue{
		Value: rdws,
		Units: extract_units(headers[17]),
	}

	plt := parse_float(row[18])
	cbcRes.PLT = cbcparser.CBCValue{
		Value: plt,
		Units: extract_units(headers[18]),
	}

	pdw := parse_float(row[19])
	cbcRes.PDW = cbcparser.CBCValue{
		Value: pdw,
		Units: extract_units(headers[19]),
	}

	mpv := parse_float(row[20])
	cbcRes.MPV = cbcparser.CBCValue{
		Value: mpv,
		Units: extract_units(headers[20]),
	}

	pct := parse_float(row[21])
	cbcRes.PCT = cbcparser.CBCValue{
		Value: pct,
		Units: extract_units(headers[21]),
	}

	plcc := parse_float(row[22])
	cbcRes.PLCC = cbcparser.CBCValue{
		Value: plcc,
		Units: extract_units(headers[22]),
	}

	plcr := parse_float(row[23])
	cbcRes.PLCR = cbcparser.CBCValue{
		Value: plcr,
		Units: extract_units(headers[23]),
	}

	// Set normal ranges if available for each CBC value
	if normal_ranges != nil {
		cbcRes.WBC.Flag = get_flag(wbc, normal_ranges.WBC)
		cbcRes.WBC.NormalRange = normal_ranges.WBC

		cbcRes.LYM.Flag = get_flag(lym, normal_ranges.LYM)
		cbcRes.LYM.NormalRange = normal_ranges.LYM

		cbcRes.LYMPercent.Flag = get_flag(lym_percent, normal_ranges.LYMPercent)
		cbcRes.LYMPercent.NormalRange = normal_ranges.LYMPercent

		cbcRes.MID.Flag = get_flag(mid, normal_ranges.MID)
		cbcRes.MID.NormalRange = normal_ranges.MID

		cbcRes.MIDPercent.Flag = get_flag(mid_percent, normal_ranges.MIDPercent)
		cbcRes.MIDPercent.NormalRange = normal_ranges.MIDPercent

		cbcRes.GRA.Flag = get_flag(granulocytes, normal_ranges.GRA)
		cbcRes.GRA.NormalRange = normal_ranges.GRA

		cbcRes.GRAPercent.Flag = get_flag(gra_percent, normal_ranges.GRAPercent)
		cbcRes.GRAPercent.NormalRange = normal_ranges.GRAPercent

		cbcRes.RBC.Flag = get_flag(rbc, normal_ranges.RBC)
		cbcRes.RBC.NormalRange = normal_ranges.RBC

		cbcRes.HGB.Flag = get_flag(hgb, normal_ranges.HGB)
		cbcRes.HGB.NormalRange = normal_ranges.HGB

		cbcRes.HCT.Flag = get_flag(hct, normal_ranges.HCT)
		cbcRes.HCT.NormalRange = normal_ranges.HCT

		cbcRes.MCV.Flag = get_flag(mcv, normal_ranges.MCV)
		cbcRes.MCV.NormalRange = normal_ranges.MCV

		cbcRes.MCH.Flag = get_flag(mch, normal_ranges.MCH)
		cbcRes.MCH.NormalRange = normal_ranges.MCH

		cbcRes.MCHC.Flag = get_flag(mchc, normal_ranges.MCHC)
		cbcRes.MCHC.NormalRange = normal_ranges.MCHC

		cbcRes.RDWc.Flag = get_flag(rdwc, normal_ranges.RDWc)
		cbcRes.RDWc.NormalRange = normal_ranges.RDWc

		cbcRes.RDWs.Flag = get_flag(rdws, normal_ranges.RDWs)
		cbcRes.RDWs.NormalRange = normal_ranges.RDWs

		cbcRes.PLT.Flag = get_flag(plt, normal_ranges.PLT)
		cbcRes.PLT.NormalRange = normal_ranges.PLT

		cbcRes.PDW.Flag = get_flag(pdw, normal_ranges.PDWc)
		cbcRes.PDW.NormalRange = normal_ranges.PDW

		cbcRes.MPV.Flag = get_flag(mpv, normal_ranges.MPV)
		cbcRes.MPV.NormalRange = normal_ranges.MPV

		cbcRes.PCT.Flag = get_flag(pct, normal_ranges.PCT)
		cbcRes.PCT.NormalRange = normal_ranges.PCT

		cbcRes.PLCC.Flag = get_flag(plcc, normal_ranges.PLCC)
		cbcRes.PLCC.NormalRange = normal_ranges.PLCC

		cbcRes.PLCR.Flag = get_flag(plcr, normal_ranges.PLCR)
		cbcRes.PLCR.NormalRange = normal_ranges.PLCR
	}

	return cbcRes
}
