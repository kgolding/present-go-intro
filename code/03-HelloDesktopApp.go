package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Hello World")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello World!"),
		widget.NewButton("Quit", func() {
			a.Quit()
		}),
	))
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}
