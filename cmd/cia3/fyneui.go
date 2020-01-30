package main

import (
	"net/url"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func fyneUi() {
	app := app.New()
	u, err := url.Parse(httpUrlString)
	if err != nil {
		errorChannel <- err
	}
	w := app.NewWindow("Civ Intelligence Agency III")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Civ Intelligence Agency III"),
		widget.NewLabel("v"+appVersion),
		widget.NewLabel("Browse to the following link"),
		widget.NewHyperlink(httpUrlString, u),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}
