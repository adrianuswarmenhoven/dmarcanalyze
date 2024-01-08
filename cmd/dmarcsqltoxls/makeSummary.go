package main

import "github.com/xuri/excelize/v2"

func makeSummary(f *excelize.File, Summaries []SheetSummary) {
	// Build the summary sheet
	loc := ""
	loc, _ = excelize.CoordinatesToCellName(2, 1)
	f.SetCellValue("data-summary", loc, "Year")
	loc, _ = excelize.CoordinatesToCellName(3, 1)
	f.SetCellValue("data-summary", loc, "Month")
	loc, _ = excelize.CoordinatesToCellName(4, 1)
	f.SetCellValue("data-summary", loc, "Count")
	loc, _ = excelize.CoordinatesToCellName(5, 1)
	f.SetCellValue("data-summary", loc, "DKIM Pass")
	loc, _ = excelize.CoordinatesToCellName(6, 1)
	f.SetCellValue("data-summary", loc, "SPF Pass")
	loc, _ = excelize.CoordinatesToCellName(7, 1)
	f.SetCellValue("data-summary", loc, "DKIM Auth Pass")
	loc, _ = excelize.CoordinatesToCellName(8, 1)
	f.SetCellValue("data-summary", loc, "SPF Auth Pass")

	lastYear := 0
	flipFlopper := false
	cellStyle := 0
	cellStyleName := ""
	for ridx, Summary := range Summaries {
		loc, _ = excelize.CoordinatesToCellName(2, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.Year)
		loc, _ = excelize.CoordinatesToCellName(3, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.Month)
		loc, _ = excelize.CoordinatesToCellName(4, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.Count)
		loc, _ = excelize.CoordinatesToCellName(5, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.DKIMPass)
		loc, _ = excelize.CoordinatesToCellName(6, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.SPFPass)
		loc, _ = excelize.CoordinatesToCellName(7, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.DKIMAuthPass)
		loc, _ = excelize.CoordinatesToCellName(8, 2+ridx)
		f.SetCellValue("data-summary", loc, Summary.SPFAuthPass)

		if lastYear != Summary.Year {
			flipFlopper = !flipFlopper
			lastYear = Summary.Year
		}
		if flipFlopper {
			if ridx%2 == 0 {
				cellStyleName = "aLight"
			} else {
				cellStyleName = "aDark"
			}
		} else {
			if ridx%2 == 0 {
				cellStyleName = "bLight"
			} else {
				cellStyleName = "bDark"
			}
		}
		cellStyle = cellStyles[cellStyleName]
		loc, _ := excelize.CoordinatesToCellName(2, 2+ridx)
		locEnd, _ := excelize.CoordinatesToCellName(8, 2+ridx)
		f.SetCellStyle("data-summary", loc, locEnd, cellStyle)
	}
	loc, _ = excelize.CoordinatesToCellName(2, 1)
	locend, _ := excelize.CoordinatesToCellName(8, 1+len(Summaries))
	f.AutoFilter("data-summary", loc+":"+locend, []excelize.AutoFilterOptions{})

	setAutoWidth(f, "data-summary")
}
