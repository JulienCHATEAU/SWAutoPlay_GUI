package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type BoolProperty struct {
	Name  string
	Value bool
}

type Dungeon struct {
	Name          string
	CreateContent func() (*gtk.Grid, error)
}

func (param *BoolProperty) ToString() string {
	return fmt.Sprintf("Name : %s - Value : %t", param.Name, param.Value)
}

var home = BoolProperty{"Phone home page", true}
var island = BoolProperty{"SW island", false}
var toa = BoolProperty{"ToA stages page", false}

var chest = BoolProperty{"Chest", true}
var sp = BoolProperty{"Social Point", false}
var crystals = BoolProperty{"Crystals", false}
var noRefill = BoolProperty{"Don't refill", false}

var dungeons = []Dungeon{
	Dungeon{"Giant", createGiantContent},
	Dungeon{"Drake", createGiantContent},
	Dungeon{"Necropolis", createGiantContent},
	Dungeon{"ToA", createGiantContent},
	Dungeon{"Scenario", createGiantContent}}

func main() {

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetDefaultSize(700, 550)
	win.SetTitle("SWAP")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	//Window Grid
	windowGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create windowGrid:", err)
	}
	windowGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	//Title
	lab, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create lab:", err)
	}
	lab.SetMarkup("<span foreground=\"#d1a432\" size=\"x-large\" face=\"serif\"><b>Summoners War Auto Play</b></span>")
	lab.SetMarginBottom(5)
	lab.SetMarginTop(5)
	windowGrid.Add(lab)

	//Tabs
	nb, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Unable to create notebook:", err)
	}
	nb.SetHExpand(true)
	nb.SetVExpand(true)
	for _, dungeon := range dungeons {
		nbChild, err := dungeon.CreateContent()
		if err != nil {
			log.Fatal("Unable to create nbChild:", err)
		}
		nbTab, err := gtk.LabelNew(dungeon.Name)
		if err != nil {
			log.Fatal("Unable to create nbTab:", err)
		}
		nb.AppendPage(nbChild, nbTab)
	}
	nb.SetMarginEnd(10)
	nb.SetMarginStart(10)
	windowGrid.Add(nb)

	//Run position grid
	props := []*BoolProperty{&home, &island, &toa}
	runPosGrid, err := createGridBoolBox("Start this run from : ", props)
	if err != nil {
		log.Fatal("createGridBoolBox() failed :", err)
	}
	windowGrid.Add(runPosGrid)

	//run button
	btn, err := gtk.ButtonNewWithLabel("Run !")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.SetMarginBottom(10)
	btn.SetMarginEnd(10)
	btn.SetMarginStart(10)
	windowGrid.Add(btn)

	win.Add(windowGrid)
	win.ShowAll()

	gtk.Main()
}

func createGridBoolBox(labelValue string, props []*BoolProperty) (*gtk.Grid, error) {
	runPosGrid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	runPosGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	runPosLabel, err := createSubTitleLabel(labelValue)
	if err != nil {
		return nil, err
	}
	runPosGrid.Add(runPosLabel)
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	var radio []*gtk.RadioButton
	radio = make([]*gtk.RadioButton, len(props))
	for index, prop := range props {
		if index == 0 {
			radio[index], _ = gtk.RadioButtonNewWithLabel(nil, prop.Name)
		} else {
			radio[index], _ = gtk.RadioButtonNewWithLabelFromWidget(radio[0], prop.Name)
		}
		radio[index].SetMarginEnd(5)
		p := prop
		r := radio[index]
		radio[index].SetActive(p.Value)
		radio[index].Connect("toggled", func() {
			updateParam(p, r.GetActive())
			log.Print(p.ToString())
		})
		box.PackStart(radio[index], true, true, 0)
	}
	runPosGrid.Add(box)
	return runPosGrid, err
}

func createGiantContent() (*gtk.Grid, error) {
	contentGrid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	contentGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	refillProps := []*BoolProperty{&chest, &sp, &crystals, &noRefill}
	refillGrid, err := createGridBoolBox("Refill energy from : ", refillProps)
	refillGrid.SetMarginTop(10)
	contentGrid.Add(refillGrid)

	return contentGrid, nil
}

func createSubTitleLabel(name string) (*gtk.Label, error) {
	label, err := gtk.LabelNew(name)
	label.SetMarginBottom(5)
	label.SetMarginTop(5)
	label.SetMarginStart(10)
	label.SetMarginEnd(25)
	return label, err
}

func updateParam(param *BoolProperty, state bool) {
	param.Value = state
}
