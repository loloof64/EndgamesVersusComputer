package main

import (
	"EndgamesVersusComputer/text"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1100,
		Height: 750,
		Title:  "Endgames versus Computer",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(text.NewTextFileManager())
	app.Run()
}
