package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dhowlett99/fixturemgr/pkg/sequence"
	"github.com/dhowlett99/list/pkg/fixture"
)

var fixtureList = []string{}
var channelList = []string{}
var fixtureInfoPanelList = []string{}
var groupList = []string{}

var selectedGroup int

const BUTTON_OUTER int = 0
const RECTANGLE int = 1
const BUTTON int = 0

func main() {
	myApp := app.New()
	w := myApp.NewWindow("DMXLights Fixture Editor")

	Blue := color.NRGBA{R: 0, G: 0, B: 100, A: 100}

	White := color.NRGBA{R: 0, G: 0, B: 0, A: 0}

	var fixturePanel *widget.List

	// Read sequences config file
	fmt.Println("Load Sequences Config File")
	sequencesConfig, err := sequence.LoadSequences()
	if err != nil {
		fmt.Printf("fixture manager: error failed to load sequences config: %s\n", err.Error())
		os.Exit(1)
	}

	// Find the group numbers.
	for _, s := range sequencesConfig.Sequences {
		group := strconv.Itoa(s.Group) //+ s.Name + s.Name + s.Type
		groupList = append(groupList, group)
	}

	fixturesConfig, err := fixture.LoadFixtures()
	if err != nil {
		fmt.Printf("fixture manager: error loading fixtures config file.\n")
		os.Exit(1)
	}

	// Setup Channels Panel.
	channelList = []string{}
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
				fmt.Printf("I have selected channel %s\n", channelList[i])
			}
		})

	// Setup Fixture Info Panel.
	fixtureInfoPanelList = []string{"Fixture", "Info"}
	fixtureInfoPanel := widget.NewList(
		func() int {
			return len(fixtureInfoPanelList)
		},
		func() fyne.CanvasObject {
			//return widget.NewLabel("template")
			return container.NewMax(
				widget.NewLabel("template"),
				widget.NewButton("", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			//o.(*widget.Label).SetText(fixtureList[i])

			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(fixtureInfoPanelList[i])

			// new part
			o.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				fmt.Println("I am selecting a " + fixtureInfoPanelList[i])
			}
		})

	// Setup Fixtures Select Panel.
	var rFixture *canvas.Rectangle
	var lastSelectedFixtureRectange *canvas.Rectangle
	fixturePanel = widget.NewList(
		func() int {
			return len(fixtureList)
		},
		func() fyne.CanvasObject {
			//return widget.NewLabel("template")
			b := container.NewMax(widget.NewButton("", nil))
			r := canvas.NewRectangle(White)
			c := container.NewMax(b, r)

			return c

		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// Configure the text on the button.
			o.(*fyne.Container).Objects[BUTTON_OUTER].(*fyne.Container).Objects[BUTTON].(*widget.Button).SetText(fixtureList[i])

			// Configure the button.
			o.(*fyne.Container).Objects[BUTTON_OUTER].(*fyne.Container).Objects[BUTTON].(*widget.Button).OnTapped = func() {
				fmt.Printf("Fixture Selection is %s\n", fixtureList[i])

				// Turn off any existing selections.
				if lastSelectedFixtureRectange != nil {
					lastSelectedFixtureRectange.FillColor = White
					lastSelectedFixtureRectange.Refresh()
				}

				// Turn on the selected button.
				rFixture = o.(*fyne.Container).Objects[RECTANGLE].(*canvas.Rectangle)
				rFixture.FillColor = Blue
				rFixture.Refresh()
				lastSelectedFixtureRectange = rFixture

				// Populate the fixtures info panel.
				fixtureInfoPanelList = getFixtureDetails(fixtureList[i], selectedGroup, fixturesConfig)
				fixtureInfoPanel.Refresh()

				// Populate the channles panel
				channelList = getChannelList(fixtureList[i], selectedGroup, fixturesConfig)
				channelPanel.Refresh()
			}
		})

	// Setup Group Select Panel.
	var rGroup *canvas.Rectangle
	var lastSelectedGroupRectange *canvas.Rectangle
	groupPanel := widget.NewList(
		func() int {
			return len(groupList)
		},
		func() fyne.CanvasObject {
			//return widget.NewLabel("template")
			b := container.NewMax(widget.NewButton("", nil))
			r := canvas.NewRectangle(White)
			c := container.NewMax(b, r)

			return c

		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// Configure the text on the button.
			o.(*fyne.Container).Objects[BUTTON_OUTER].(*fyne.Container).Objects[BUTTON].(*widget.Button).SetText(groupList[i])

			// Configure the button.
			o.(*fyne.Container).Objects[BUTTON_OUTER].(*fyne.Container).Objects[BUTTON].(*widget.Button).OnTapped = func() {
				fmt.Printf("I haves selected group no %s\n", groupList[i])

				selectedGroup, _ = strconv.Atoi(groupList[i])

				// Turn off any existing selections.
				if lastSelectedGroupRectange != nil {
					lastSelectedGroupRectange.FillColor = White
					lastSelectedGroupRectange.Refresh()
				}

				// Turn on the selected button.
				rGroup = o.(*fyne.Container).Objects[RECTANGLE].(*canvas.Rectangle)
				rGroup.FillColor = Blue
				rGroup.Refresh()
				lastSelectedGroupRectange = rGroup

				// Populate the fixtures list panel based on the group number
				fixtureList = getFixtureList(groupList[i], fixturesConfig)
				fixturePanel.Refresh()

				// Populate the fixtures info based on the first fixture in this group.
				fixtureInfoPanelList = getFixtureDetails("1", selectedGroup, fixturesConfig)
				fixtureInfoPanel.Refresh()
			}
		})

	w.Resize(fyne.NewSize(800, 400))
	panel1 := container.NewHSplit(fixtureInfoPanel, channelPanel)
	panel1.Offset = 0.5
	panel2 := container.NewHSplit(fixturePanel, panel1)
	panel2.Offset = 0.1
	mainPanel := container.NewHSplit(groupPanel, panel2)
	mainPanel.Offset = 0.1

	w.SetContent(mainPanel)
	w.ShowAndRun()
}

func getFixtureList(groupNumber string, fixturesConfig *fixture.Fixtures) []string {

	group, _ := strconv.Atoi(groupNumber)

	fixtureList = []string{}
	for _, f := range fixturesConfig.Fixtures {

		if f.Group == group {
			fixtureList = append(fixtureList, fmt.Sprintf("%d", f.Number))
		}

	}
	return fixtureList
}

func getChannelList(fixtureNumber string, groupNumber int, fixturesConfig *fixture.Fixtures) []string {

	fixture, _ := strconv.Atoi(fixtureNumber)

	channelsList := []string{}

	fmt.Printf("Fixture Number %d Group Number %d\n", fixture, groupNumber)

	for _, f := range fixturesConfig.Fixtures {
		//fmt.Printf("Fixture No %d %s  %s \n", f.Number, f.Name, f.Description)

		if f.Group == groupNumber {
			if f.Number == fixture {
				for _, c := range f.Channels {
					//fmt.Printf("Channel %+v\n", c.Name)
					channelsList = append(channelsList, c.Name)
				}
			}
		}

	}
	return channelsList
}

func getFixtureDetails(number string, group int, fixturesConfig *fixture.Fixtures) []string {

	fixtureNumber, _ := strconv.Atoi(number)
	fixtureInfoPanelList := []string{}

	fmt.Printf("Get Fixture Details for Fixture Number %d Group Number %d \n", fixtureNumber, group)

	for _, f := range fixturesConfig.Fixtures {
		//fmt.Printf("Fixture No %d %s  %s \n", f.Number, f.Name, f.Description)

		if f.Group == group {
			if f.Number == fixtureNumber {
				fixtureInfoPanelList = []string{"Fixture Number:" + number, "Name:" + f.Name, "Description:" + f.Description, "Group:" + strconv.Itoa(f.Group)}
			}
		}
	}
	return fixtureInfoPanelList
}
