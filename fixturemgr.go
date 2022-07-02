package main

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowlett99/list/pkg/fixture"
)

var fixtureList = []string{}
var channelList = []string{}

const BUTTON_OUTER int = 1
const RECTANGLE int = 0
const BUTTON int = 0

func main() {
	myApp := app.New()
	w := myApp.NewWindow("DMXLights Fixture Editor")

	Blue := color.NRGBA{R: 0, G: 0, B: 100, A: 100}

	White := color.NRGBA{R: 0, G: 0, B: 0, A: 0}

	content := container.NewMax()

	var fixturePanel *widget.List

	fixtureConfig, err := fixture.LoadFixtures()
	if err != nil {
		fmt.Printf("error loading fixtures.\n")
		os.Exit(1)
	}

	for _, f := range fixtureConfig.Fixtures {
		if f.Type == "rgb" {
			fixtureList = append(fixtureList, fmt.Sprintf("%d", f.Number))

		}
	}

	// Setup Channels Panel.
	channelList = []string{"one", "two"}
	channelPanel := widget.NewList(
		func() int {
			return len(channelList)
		},
		func() fyne.CanvasObject {
			//return widget.NewLabel("template")
			return container.NewPadded(
				widget.NewLabel("template"),
				widget.NewButton("", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			//o.(*widget.Label).SetText(fixtureList[i])

			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(channelList[i])

			// new part
			o.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				fmt.Println("I am button " + channelList[i])
			}
		})

	// Setup Fixtures Panel.
	var r *canvas.Rectangle
	var lastSelectedRectange *canvas.Rectangle
	fixturePanel = widget.NewList(
		func() int {
			return len(fixtureList)
		},
		func() fyne.CanvasObject {
			//return widget.NewLabel("template")
			b := container.NewPadded(widget.NewButton("", nil))
			r := canvas.NewRectangle(White)
			c := container.NewMax(r, b)

			return c

		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// Configure the text on the button.
			o.(*fyne.Container).Objects[BUTTON_OUTER].(*fyne.Container).Objects[BUTTON].(*widget.Button).SetText(fixtureList[i])

			// Configure the button.
			o.(*fyne.Container).Objects[BUTTON_OUTER].(*fyne.Container).Objects[BUTTON].(*widget.Button).OnTapped = func() {
				fmt.Println("I am button " + fixtureList[i])

				// Turn off any existing selections.
				if lastSelectedRectange != nil {
					lastSelectedRectange.FillColor = White
					lastSelectedRectange.Refresh()
				}

				// Turn on the selected button.
				r = o.(*fyne.Container).Objects[RECTANGLE].(*canvas.Rectangle)
				r.FillColor = Blue
				r.Refresh()
				lastSelectedRectange = r

				channelList = getChannelList(fixtureList[i])
				channelPanel.Refresh()
			}
		})

	w.Resize(fyne.NewSize(800, 500))
	nextPanel := container.NewHSplit(channelPanel, content)
	mainPanel := container.NewHSplit(fixturePanel, nextPanel)
	mainPanel.Offset = 0.33

	w.SetContent(mainPanel)
	w.ShowAndRun()
}

func getChannelList(fixtureNumber string) []string {

	out := []string{}

	out = append(out, fmt.Sprintf("Number %s", fixtureNumber))
	return out
}
