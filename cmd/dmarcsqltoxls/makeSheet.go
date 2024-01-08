package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

var (
	cellStyles = map[string]int{}
	header     = []string{
		"Report ID",
		"Org Name",
		"Email",
		"Extra Contact Info",
		"Begin Date/Time",
		"End Date/Time",
		"Policy Published Domain",
		"Policy Published ADKIM",
		"Policy Published ASPF",
		"Policy Published Percentage",
		"Policy Published Policy",
		"Policy Published Subdomain Policy",
		"Source IP",
		"Count",
		"Disposition",
		"DKIM",
		"SPF",
		"Header From",
		"DKIM Auth Result Domain",
		"DKIM Auth Result Result",
		"DKIM Auth Result Selector",
		"SPF Auth Result Domain",
		"SPF Auth Result Result",
		"SPF Auth Result Scope",
	}
)

type SheetSummary struct {
	Year         int
	Month        int
	Rows         int
	Count        int
	DKIMPass     int
	SPFPass      int
	DKIMAuthPass int
	SPFAuthPass  int
}

func makeSheet(f *excelize.File,
	year int, month int,
	TimeBasedIndex map[int]map[int][]string,
	MetaDataIndex map[string]*Metadata,
	PolicyPublishedIndex map[string]*PolicyPublished,
	RecordIndex map[string][]*Record) SheetSummary {

	Summary := SheetSummary{
		Year:  year,
		Month: month,
		Rows:  len(TimeBasedIndex[year][month]),
	}

	// Create raw data for Raw Data sheet
	// Group by ReportID
	rawSheet := make([][]interface{}, 0)
	sheetName := fmt.Sprintf("%d-%02d", year, month)
	for _, reportID := range TimeBasedIndex[year][month] {
		m := MetaDataIndex[reportID]

		rowMetaData := make([]interface{}, 0)
		rowMetaData = append(rowMetaData, m.ReportID)
		rowMetaData = append(rowMetaData, m.OrgName)
		rowMetaData = append(rowMetaData, m.Email)
		rowMetaData = append(rowMetaData, m.ExtraContactInfo)
		rowMetaData = append(rowMetaData, time.Unix(m.Begin, 0))
		rowMetaData = append(rowMetaData, time.Unix(m.End, 0))
		p := PolicyPublishedIndex[m.ReportID]
		rowMetaData = append(rowMetaData, p.Domain)
		rowMetaData = append(rowMetaData, p.ADKIM)
		rowMetaData = append(rowMetaData, p.ASPF)
		rowMetaData = append(rowMetaData, p.Percentage)
		rowMetaData = append(rowMetaData, p.Policy)
		rowMetaData = append(rowMetaData, p.SPolicy)

		// Now copy that base row for each record and add the record data
		for _, r := range RecordIndex[m.ReportID] {
			row := rowMetaData[:] // Clone the base
			row = append(row, r.SourceIP)
			row = append(row, r.Count)
			row = append(row, r.Disposition)
			row = append(row, r.DKIM)
			row = append(row, r.SPF)
			row = append(row, r.HeaderFrom)
			row = append(row, r.DKIMAuthResultDomain)
			row = append(row, r.DKIMAuthResultResult)
			row = append(row, r.DKIMAuthResultSelector)
			row = append(row, r.SPFAuthResultDomain)
			row = append(row, r.SPFAuthResultResult)
			row = append(row, r.SPFAuthResultScope)

			rawSheet = append(rawSheet, row)

			// Update summary
			Summary.Count += r.Count
			if r.DKIM == "pass" {
				Summary.DKIMPass += r.Count
			}
			if r.SPF == "pass" {
				Summary.SPFPass += r.Count
			}
			if r.DKIMAuthResultResult == "pass" {
				Summary.DKIMAuthPass += r.Count
			}
			if r.SPFAuthResultResult == "pass" {
				Summary.SPFAuthPass += r.Count
			}
		}
	}

	// Create a new sheet.
	_, err := f.NewSheet(sheetName)
	if err != nil {
		slog.Error("error creating sheet", "error", err)
		os.Exit(1)
	}

	loc, _ := excelize.CoordinatesToCellName(2, 1)
	f.SetSheetRow(sheetName, loc, &header)
	// Set value of a row.
	reportID := ""
	flipFlopper := true

	cellStyleName := ""
	cellStyle := 0
	cellDateStyleName := ""
	cellDateStyle := 0

	for ridx, row := range rawSheet {
		loc, _ := excelize.CoordinatesToCellName(2, 2+ridx)
		locEnd, _ := excelize.CoordinatesToCellName(2+len(header), 2+ridx)
		locDateStart, _ := excelize.CoordinatesToCellName(6, 2+ridx)
		locDateEnd, _ := excelize.CoordinatesToCellName(7, 2+ridx)
		if reportID != row[0].(string) {
			flipFlopper = !flipFlopper
			reportID = row[0].(string)
		}
		if flipFlopper {
			if ridx%2 == 0 {
				cellStyleName = "aLight"
				cellDateStyleName = "aLightDate"
			} else {
				cellStyleName = "aDark"
				cellDateStyleName = "aDarkDate"
			}
		} else {
			if ridx%2 == 0 {
				cellStyleName = "bLight"
				cellDateStyleName = "bLightDate"
			} else {
				cellStyleName = "bDark"
				cellDateStyleName = "bDarkDate"
			}
		}

		if !(row[15].(string) == "pass" && row[16].(string) == "pass") {
			cellStyleName += "Fail"
			cellDateStyleName += "Fail"
		}
		cellStyle = cellStyles[cellStyleName]
		cellDateStyle = cellStyles[cellDateStyleName]

		f.SetCellStyle(sheetName, loc, locEnd, cellStyle)
		f.SetSheetRow(sheetName, loc, &row)
		f.SetCellStyle(sheetName, locDateStart, locDateEnd, cellDateStyle)

	}

	setAutoWidth(f, sheetName)

	// Make a table out of the raw data
	loc, _ = excelize.CoordinatesToCellName(2, 1)
	locend, _ := excelize.CoordinatesToCellName(2+len(header), 1+len(rawSheet))
	f.AutoFilter(sheetName, loc+":"+locend, []excelize.AutoFilterOptions{})

	return Summary
}

func setAutoWidth(f *excelize.File, sheetName string) {
	// Autofit all columns according to their text content
	cols, err := f.GetCols(sheetName)
	if err != nil {
		return
	}
	for idx, col := range cols {
		if idx == 0 {
			continue
		}
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			slog.Error("error converting column number to name", "error", err)
			os.Exit(1)
		}
		f.SetColWidth(sheetName, name, name, float64(largestWidth))
	}
}
