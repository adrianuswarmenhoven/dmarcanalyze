package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/xuri/excelize/v2"
)

func buildXLS() {
	slog.Info("Fetching data from database")
	MetaDatas, err := db.FetchMetadata()
	if err != nil {
		slog.Error("error fetching metadata", "error", err)
		os.Exit(1)
	}
	PoliciesPublished, err := db.FetchPolicyPublished()
	if err != nil {
		slog.Error("error fetching policy published", "error", err)
		os.Exit(1)
	}
	Records, err := db.FetchRecords()
	if err != nil {
		slog.Error("error fetching records", "error", err)
		os.Exit(1)
	}
	err = db.Close()
	if err != nil {
		slog.Error("error closing database", "error", err)
		os.Exit(1)
	}
	// Done with DB, first make some indices

	// YYYY -> MM -> ReportID
	TimeBasedIndex := make(map[int]map[int][]string)

	MetaDataIndex := make(map[string]*Metadata)
	for _, m := range MetaDatas {
		year := time.Unix(m.Begin, 0).Year()
		month := int(time.Unix(m.Begin, 0).Month())
		if _, ok := TimeBasedIndex[year]; !ok {
			TimeBasedIndex[year] = make(map[int][]string)
		}
		if _, ok := TimeBasedIndex[year][month]; !ok {
			TimeBasedIndex[year][month] = make([]string, 0)
		}
		TimeBasedIndex[year][month] = append(TimeBasedIndex[year][month], m.ReportID)
		MetaDataIndex[m.ReportID] = m
	}
	PolicyPublishedIndex := make(map[string]*PolicyPublished)
	for _, p := range PoliciesPublished {
		PolicyPublishedIndex[p.ReportID] = p
	}
	RecordIndex := make(map[string][]*Record)
	for _, r := range Records {
		RecordIndex[r.ReportID] = append(RecordIndex[r.ReportID], r)
	}

	// Now we have indices, we can build the XLSX
	YearIndex := make([]int, 0)
	for year, _ := range TimeBasedIndex {
		YearIndex = append(YearIndex, year)
	}
	// Sort the years
	slices.Sort(YearIndex)
	slices.Reverse(YearIndex)
	slog.Info("Building XLSX")

	var f *excelize.File
	if Configuration.XLS.Template == "" {
		slog.Info("Using blank template")
		f = excelize.NewFile()
	} else {
		slog.Info("Using template", "template", Configuration.XLS.Template)
		f, err = excelize.OpenFile(Configuration.XLS.Template)
		if err != nil {
			slog.Error("error opening template", "error", err)
			os.Exit(1)
		}
	}

	setSheetStyle(f)
	defer func() {
		slog.Info("Saving XLSX", "filename", Configuration.XLS.Output)
		if err := f.Close(); err != nil {
			slog.Error("error closing file", "error", err)
		}
	}()

	f.DeleteSheet("data-summary") // Remove any placeholders for the template
	firstsheet, err := f.NewSheet("data-summary")
	if err != nil {
		slog.Error("error creating sheet", "error", err)
		os.Exit(1)
	}
	// Create a new bookend sheet for the first data sheet
	// This is so we can use generic formula's in the template
	_, err = f.NewSheet("data-first")
	if err != nil {
		slog.Error("error creating sheet", "error", err)
		os.Exit(1)
	}
	Summaries := make([]SheetSummary, 0)
	for idx, year := range YearIndex {
		slog.Debug("Building sheet", "progress", fmt.Sprintf("%.2f%%", float64(idx)/float64(len(YearIndex))*100))
		MonthIndex := make([]int, 0)
		for month, _ := range TimeBasedIndex[year] {
			MonthIndex = append(MonthIndex, month)
		}
		// Sort the months
		slices.Sort(MonthIndex)
		slices.Reverse(MonthIndex)
		for _, month := range MonthIndex {
			Summary := makeSheet(f, year, month, TimeBasedIndex, MetaDataIndex, PolicyPublishedIndex, RecordIndex)
			Summaries = append(Summaries, Summary)
		}
	}
	// Create a new bookend sheet for the first data sheet
	// This is so we can use generic formula's in the template
	_, err = f.NewSheet("data-last")
	if err != nil {
		slog.Error("error creating sheet", "error", err)
		os.Exit(1)
	}

	slog.Info("Building summary")
	makeSummary(f, Summaries)

	f.SetActiveSheet(firstsheet)

	if err := f.SaveAs(Configuration.XLS.Output); err != nil {
		slog.Error("error saving file", "error", err)
		os.Exit(1)
	}

}
