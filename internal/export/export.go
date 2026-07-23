// Package export renders ServiceNow records (as returned by internal/client)
// to JSON, CSV, XML, or XLSX.
package export

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
	FormatXML  Format = "xml"
	FormatXLSX Format = "xlsx"
)

// ParseFormat validates a user-supplied format string.
func ParseFormat(s string) (Format, error) {
	switch f := Format(strings.ToLower(s)); f {
	case FormatJSON, FormatCSV, FormatXML, FormatXLSX:
		return f, nil
	default:
		return "", fmt.Errorf("unsupported format %q (want one of: json, csv, xml, xlsx)", s)
	}
}

// WriteRecords renders records to w in the given format. For the tabular
// formats (csv, xml, xlsx), columns are the union of all record keys,
// sorted for stable, repeatable output across runs.
func WriteRecords(w io.Writer, format Format, records []map[string]any) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, records)
	case FormatCSV:
		return writeCSV(w, records)
	case FormatXML:
		return writeXML(w, records)
	case FormatXLSX:
		return writeXLSX(w, records)
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func writeJSON(w io.Writer, records []map[string]any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func writeCSV(w io.Writer, records []map[string]any) error {
	fields := columns(records)

	cw := csv.NewWriter(w)
	if err := cw.Write(fields); err != nil {
		return err
	}
	for _, rec := range records {
		row := make([]string, len(fields))
		for i, f := range fields {
			row[i] = stringify(rec[f])
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeXML(w io.Writer, records []map[string]any) error {
	fields := columns(records)

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")

	root := xml.StartElement{Name: xml.Name{Local: "result"}}
	if err := enc.EncodeToken(root); err != nil {
		return err
	}
	for _, rec := range records {
		recElem := xml.StartElement{Name: xml.Name{Local: "record"}}
		if err := enc.EncodeToken(recElem); err != nil {
			return err
		}
		for _, f := range fields {
			fieldElem := xml.StartElement{Name: xml.Name{Local: xmlName(f)}}
			if err := enc.EncodeToken(fieldElem); err != nil {
				return err
			}
			if err := enc.EncodeToken(xml.CharData([]byte(stringify(rec[f])))); err != nil {
				return err
			}
			if err := enc.EncodeToken(fieldElem.End()); err != nil {
				return err
			}
		}
		if err := enc.EncodeToken(recElem.End()); err != nil {
			return err
		}
	}
	if err := enc.EncodeToken(root.End()); err != nil {
		return err
	}
	return enc.Flush()
}

func writeXLSX(w io.Writer, records []map[string]any) error {
	fields := columns(records)

	f := excelize.NewFile()
	defer f.Close()

	const sheet = "Records"
	f.SetSheetName(f.GetSheetName(0), sheet)

	for i, field := range fields {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return err
		}
		if err := f.SetCellValue(sheet, cell, field); err != nil {
			return err
		}
	}
	for r, rec := range records {
		for c, field := range fields {
			cell, err := excelize.CoordinatesToCellName(c+1, r+2)
			if err != nil {
				return err
			}
			if err := f.SetCellValue(sheet, cell, stringify(rec[field])); err != nil {
				return err
			}
		}
	}

	return f.Write(w)
}

// columns returns the union of all keys across records, sorted alphabetically.
// JSON decoding into map[string]any does not preserve field order, so a
// deterministic order is used instead of an arbitrary "first seen" one.
func columns(records []map[string]any) []string {
	set := make(map[string]struct{})
	for _, rec := range records {
		for k := range rec {
			set[k] = struct{}{}
		}
	}
	cols := make([]string, 0, len(set))
	for k := range set {
		cols = append(cols, k)
	}
	sort.Strings(cols)
	return cols
}

// stringify flattens a decoded JSON value into a single cell/element value.
// ServiceNow reference fields decode as {"link": "...", "value": "sys_id"};
// the sys_id is used since it is the stable, meaningful part for a flat cell.
func stringify(v any) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case map[string]any:
		if s, ok := val["value"].(string); ok {
			return s
		}
		b, _ := json.Marshal(val)
		return string(b)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprint(val)
		}
		return string(b)
	}
}

// xmlName guards against field names that aren't valid XML element names
// (e.g. empty, or starting with a digit).
func xmlName(field string) string {
	if field == "" {
		return "_"
	}
	if field[0] >= '0' && field[0] <= '9' {
		return "_" + field
	}
	return field
}
