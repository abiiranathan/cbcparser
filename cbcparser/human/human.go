package human

// parses cbc tsv/csv to print the cbc report to stdout
// Ascii char go from:
// https://theasciicode.com.ar/extended-ascii-code/box-drawing-character-single-line-upper-left-corner-ascii-code-218.html
//
// CBC format:
// Sample ID	Date	Time	Patient ID	Birth date	WBC 10^9/l	WBC flag	LYM 10^9/l	LYM flag	MID 10^9/l	MID flag	GRA 10^9/l	GRA flag	LYM% %	LYM% flag	MID% %	MID% flag	GRA% %	GRA% flag	RBC 10^12/l	RBC flag	HGB g/dl	HGB flag	HCT %	HCT flag	MCV fl	MCV flag	MCH pg	MCH flag	MCHC g/dl	MCHC flag	RDWs fl	RDWs flag	RDWc %	RDWc flag	PLT 10^9/l	PLT flag	PCT %	PCT flag	MPV fl	MPV flag	PDWs fl	PDWs flag	PDWc %	PDWc flag	P-LCC 10^9/l	P-LCC flag	P-LCR %	P-LCR flag	Type	Warning
//

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/abiiranathan/cbcparser/cbcparser"
)

const (
	nfields   = 52
	separator = '\t'
)

func parse_float(value string) float32 {
	fvalue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0.0
	}
	return float32(fvalue)
}

func extract_units(value string) string {
	valarr := strings.Split(value, " ")
	if len(valarr) != 2 {
		return ""
	}
	return valarr[1]
}

func New() cbcparser.CSVParser {
	return &humanCBCResult{}
}

func NewMultiParser() cbcparser.CSVMultiParser {
	return &humanCBCResultMulti{}
}

func (cbc humanCBCResult) Write(out io.Writer, format cbcparser.OutFormat) error {
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

// Parse reads from r and parses the data into an slice of a CBCWriter struct.
// The text file is expected to be in the tab-separated format.
// The first line of the file is expected to be the header.
//
// The header is expected to be in the following format(with 52 columns):
//
// Sample ID	Date	Time	Patient ID	Birth date	WBC 10^9/l	WBC flag	LYM 10^9/l	LYM flag	MID 10^9/l	MID flag	GRA 10^9/l	GRA flag	LYM% %	LYM% flag	MID% %	MID% flag	GRA% %	GRA% flag	RBC 10^12/l	RBC flag	HGB g/dl	HGB flag	HCT %	HCT flag	MCV fl	MCV flag	MCH pg	MCH flag	MCHC g/dl	MCHC flag	RDWs fl	RDWs flag	RDWc %	RDWc flag	PLT 10^9/l	PLT flag	PCT %	PCT flag	MPV fl	MPV flag	PDWs fl	PDWs flag	PDWc %	PDWc flag	P-LCC 10^9/l	P-LCC flag	P-LCR %	P-LCR flag	Type	Warning
func (cbc humanCBCResult) Parse(r io.Reader, normal_ranges *cbcparser.CBCNormalRange) (cbcparser.CBCWriter, error) {
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

	if row[0] == "0" || row[len(row)-2] == "Blank" {
		return nil, cbcparser.ErrBlankCBCRecord
	}

	result := set_cbc_value(headers, row, normal_ranges)
	return result, nil
}

// MultiParse reads from r and parses the data into an slice of a CBCWriter struct.
func (humanCBCResultMulti) ParseMulti(r io.Reader, normal_ranges *cbcparser.CBCNormalRange) (cbcparser.CBCMultiWriter, error) {
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

	for _, row := range data[1:] {
		if row[0] == "0" || row[len(row)-2] == "Blank" {
			continue
		}

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
	cbcRes := humanCBCResult{}
	// Patient Identifiers
	cbcRes.SampleID = row[0]
	cbcRes.Date = row[1]
	cbcRes.Time = row[2]
	cbcRes.PatientID = row[3]
	cbcRes.BirthDate = row[4]

	// Absoluet counts
	cbcRes.WBC = cbcparser.CBCValue{
		Value: parse_float(row[5]),
		Units: extract_units(headers[5]),
		Flag:  strings.TrimSpace(row[6]),
	}

	cbcRes.LYM = cbcparser.CBCValue{
		Value: parse_float(row[7]),
		Units: extract_units(headers[7]),
		Flag:  strings.TrimSpace(row[8]),
	}

	cbcRes.MID = cbcparser.CBCValue{
		Value: parse_float(row[9]),
		Units: extract_units(headers[9]),
		Flag:  strings.TrimSpace(row[10]),
	}

	cbcRes.GRA = cbcparser.CBCValue{
		Value: parse_float(row[11]),
		Units: extract_units(headers[11]),
		Flag:  strings.TrimSpace(row[12]),
	}

	// Percentages
	cbcRes.LYMPercent = cbcparser.CBCValue{
		Value: parse_float(row[13]),
		Units: extract_units(headers[13]),
		Flag:  strings.TrimSpace(row[14]),
	}

	cbcRes.MIDPercent = cbcparser.CBCValue{
		Value: parse_float(row[15]),
		Units: extract_units(headers[15]),
		Flag:  strings.TrimSpace(row[16]),
	}

	cbcRes.GRAPercent = cbcparser.CBCValue{
		Value: parse_float(row[17]),
		Units: extract_units(headers[17]),
		Flag:  strings.TrimSpace(row[18]),
	}

	// RBCS
	cbcRes.RBC = cbcparser.CBCValue{
		Value: parse_float(row[19]),
		Units: extract_units(headers[19]),
		Flag:  strings.TrimSpace(row[20]),
	}

	cbcRes.HGB = cbcparser.CBCValue{
		Value: parse_float(row[21]),
		Units: extract_units(headers[21]),
		Flag:  strings.TrimSpace(row[22]),
	}

	cbcRes.HCT = cbcparser.CBCValue{
		Value: parse_float(row[23]),
		Units: extract_units(headers[23]),
		Flag:  strings.TrimSpace(row[24]),
	}

	cbcRes.MCV = cbcparser.CBCValue{
		Value: parse_float(row[25]),
		Units: extract_units(headers[25]),
		Flag:  strings.TrimSpace(row[26]),
	}

	cbcRes.MCH = cbcparser.CBCValue{
		Value: parse_float(row[27]),
		Units: extract_units(headers[27]),
		Flag:  strings.TrimSpace(row[28]),
	}

	cbcRes.MCHC = cbcparser.CBCValue{
		Value: parse_float(row[29]),
		Units: extract_units(headers[29]),
		Flag:  strings.TrimSpace(row[30]),
	}

	cbcRes.RDWs = cbcparser.CBCValue{
		Value: parse_float(row[31]),
		Units: extract_units(headers[31]),
		Flag:  strings.TrimSpace(row[32]),
	}

	cbcRes.RDWc = cbcparser.CBCValue{
		Value: parse_float(row[33]),
		Units: extract_units(headers[33]),
		Flag:  strings.TrimSpace(row[34]),
	}

	cbcRes.PLT = cbcparser.CBCValue{
		Value: parse_float(row[35]),
		Units: extract_units(headers[35]),
		Flag:  strings.TrimSpace(row[36]),
	}

	cbcRes.PCT = cbcparser.CBCValue{
		Value: parse_float(row[37]),
		Units: extract_units(headers[37]),
		Flag:  strings.TrimSpace(row[38]),
	}

	cbcRes.MPV = cbcparser.CBCValue{
		Value: parse_float(row[39]),
		Units: extract_units(headers[39]),
		Flag:  strings.TrimSpace(row[40]),
	}

	cbcRes.PDWs = cbcparser.CBCValue{
		Value: parse_float(row[41]),
		Units: extract_units(headers[41]),
		Flag:  strings.TrimSpace(row[42]),
	}

	cbcRes.PDWc = cbcparser.CBCValue{
		Value: parse_float(row[43]),
		Units: extract_units(headers[43]),
		Flag:  strings.TrimSpace(row[44]),
	}

	cbcRes.PLCC = cbcparser.CBCValue{
		Value: parse_float(row[45]),
		Units: extract_units(headers[45]),
		Flag:  strings.TrimSpace(row[46]),
	}
	cbcRes.PLCR = cbcparser.CBCValue{
		Value: parse_float(row[47]),
		Units: extract_units(headers[47]),
		Flag:  strings.TrimSpace(row[48]),
	}

	cbcRes.Type = row[49]
	cbcRes.Warning = row[50]

	if normal_ranges != nil {
		cbcRes.WBC.NormalRange = normal_ranges.WBC
		cbcRes.LYM.NormalRange = normal_ranges.LYM
		cbcRes.LYMPercent.NormalRange = normal_ranges.LYMPercent
		cbcRes.MID.NormalRange = normal_ranges.MID
		cbcRes.MIDPercent.NormalRange = normal_ranges.MIDPercent
		cbcRes.GRA.NormalRange = normal_ranges.GRA
		cbcRes.GRAPercent.NormalRange = normal_ranges.GRAPercent
		cbcRes.RBC.NormalRange = normal_ranges.RBC
		cbcRes.HGB.NormalRange = normal_ranges.HGB
		cbcRes.HCT.NormalRange = normal_ranges.HCT
		cbcRes.MCV.NormalRange = normal_ranges.MCV
		cbcRes.MCH.NormalRange = normal_ranges.MCH
		cbcRes.MCHC.NormalRange = normal_ranges.MCHC
		cbcRes.RDWc.NormalRange = normal_ranges.RDWc
		cbcRes.RDWs.NormalRange = normal_ranges.RDWs
		cbcRes.PLT.NormalRange = normal_ranges.PLT
		cbcRes.PDWs.NormalRange = normal_ranges.PDWs
		cbcRes.PDWc.NormalRange = normal_ranges.PDWc
		cbcRes.MPV.NormalRange = normal_ranges.MPV
		cbcRes.PCT.NormalRange = normal_ranges.PCT
		cbcRes.PLCC.NormalRange = normal_ranges.PLCC
		cbcRes.PLCR.NormalRange = normal_ranges.PLCR
	}

	return cbcRes
}
