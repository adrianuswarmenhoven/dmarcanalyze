package main

import "github.com/xuri/excelize/v2"

func setSheetStyle(f *excelize.File) {
	// A
	// Light
	aLight, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundLight},
			Pattern: 1,
		},
	})
	cellStyles["aLight"] = aLight

	aLightFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundLight},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowAFailFontColor,
		},
	})
	cellStyles["aLightFail"] = aLightFail

	aLightDate, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundLight},
			Pattern: 1,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["aLightDate"] = aLightDate

	aLightDateFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundLight},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowAFailFontColor,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["aLightDateFail"] = aLightDateFail

	// Dark
	aDark, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundDark},
			Pattern: 1,
		},
	})
	cellStyles["aDark"] = aDark

	aDarkFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundDark},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowAFailFontColor,
		},
	})
	cellStyles["aDarkFail"] = aDarkFail

	aDarkDate, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundDark},
			Pattern: 1,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["aDarkDate"] = aDarkDate

	aDarkDateFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowABackgroundDark},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowAFailFontColor,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["aDarkDateFail"] = aDarkDateFail

	// B
	// Light
	bLight, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundLight},
			Pattern: 1,
		},
	})
	cellStyles["bLight"] = bLight

	bLightFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundLight},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowBFailFontColor,
		},
	})
	cellStyles["bLightFail"] = bLightFail

	bLightDate, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundLight},
			Pattern: 1,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["bLightDate"] = bLightDate

	bLightDateFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundLight},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowBFailFontColor,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["bLightDateFail"] = bLightDateFail

	// Dark
	bDark, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundDark},
			Pattern: 1,
		},
	})
	cellStyles["bDark"] = bDark

	bDarkFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundDark},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowBFailFontColor,
		},
	})
	cellStyles["bDarkFail"] = bDarkFail

	bDarkDate, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundDark},
			Pattern: 1,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["bDarkDate"] = bDarkDate

	bDarkDateFail, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{Configuration.XLS.Style.Datasheet.RowBBackgroundDark},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Color: Configuration.XLS.Style.Datasheet.RowBFailFontColor,
		},
		CustomNumFmt: &Configuration.XLS.Style.Datasheet.DateFormat,
	})
	cellStyles["bDarkDateFail"] = bDarkDateFail

}
