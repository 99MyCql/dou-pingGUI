package main

import (
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"

	"github.com/99MyCql/dou-pingGUI/backend"
)

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    600,
		Title:     "dou-pingGUI",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})
	app.Bind(backend.NewController())
	if err := app.Run(); err != nil {
		panic(err)
	}
}
