package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed app.theme.json
var themeData string

func loadTheme(a fyne.App) {
	if themeData == "" {
		return
	}

	th, err := theme.FromJSON(themeData)
	if err != nil {
		fyne.LogError("Failed to parse app.theme.json", err)
		return
	}

	a.Settings().SetTheme(th)
}
