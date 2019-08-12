package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/gotk3/gotk3/gtk"
)

type BoolProperty struct {
	Name        string
	Value       bool
	StringValue string
}

func (param *BoolProperty) toString() string {
	return fmt.Sprintf("Name : %s - Value : %t", param.Name, param.Value)
}

type Dungeon struct {
	Name           string
	ConcernedParam [5]bool //AverageDungeonTime | RunCount | Refill | Difficulty | StartStage
}

var normal = BoolProperty{"Normal", true, "Normal"}
var hard = BoolProperty{"Hard", false, "Hard"}
var hell = BoolProperty{"Hell", false, "Hell"}

var home = BoolProperty{"Phone home page", true, "Home"}
var island = BoolProperty{"Island", false, "Island"}
var toa = BoolProperty{"ToA stages page", false, "ToA"}

var chest = BoolProperty{"Chest", true, "Chest"}
var sp = BoolProperty{"Social Point", false, "SocialPoint"}
var crystals = BoolProperty{"Crystals", false, "Crystals"}
var noRefill = BoolProperty{"Don't refill", false, "Off"}

var props = []*BoolProperty{&home, &island, &toa}
var refillProps = []*BoolProperty{&chest, &sp, &crystals, &noRefill}
var difficultyProps = []*BoolProperty{&normal, &hard, &hell}

var dungeons = []Dungeon{
	Dungeon{"Giant", [5]bool{true, true, true, false, false}},
	Dungeon{"Drake", [5]bool{true, true, true, false, false}},
	Dungeon{"Necropolis", [5]bool{true, true, true, false, false}},
	Dungeon{"ToA", [5]bool{true, true, true, true, true}},
	Dungeon{"Scenario", [5]bool{true, true, true, true, true}},
}

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

	dungeonsTabs, _ := gtk.NotebookNew()

	dunLength := len(dungeons)
	var adts = make([]*gtk.Entry, dunLength)
	var runCounts = make([]*gtk.Entry, dunLength)
	var startStages = make([]*gtk.Entry, dunLength)

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
	if err != nil {
		log.Fatal("Unable to create notebook:", err)
	}
	dungeonsTabs.SetHExpand(true)
	dungeonsTabs.SetVExpand(true)

	for count, dungeon := range dungeons {
		contentGrid, _ := dungeon.createDungeonContent(count, adts, runCounts, startStages)

		dungeonsTabsTab, err := gtk.LabelNew(dungeon.Name)
		if err != nil {
			log.Fatal("Unable to create dungeonsTabsTab:", err)
		}
		dungeonsTabs.AppendPage(contentGrid, dungeonsTabsTab)
	}

	dungeonsTabs.SetMarginEnd(10)
	dungeonsTabs.SetMarginStart(10)
	windowGrid.Add(dungeonsTabs)

	//Run position grid
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
	btn.Connect("clicked", func() {
		run(dungeonsTabs, adts, runCounts, startStages)
	})
	windowGrid.Add(btn)

	win.Add(windowGrid)
	win.ShowAll()

	gtk.Main()
}

func run(dungeonsTabs *gtk.Notebook, adt []*gtk.Entry, runCount []*gtk.Entry, startStage []*gtk.Entry) { //AverageDungeonTime | RunCount | Refill | Difficulty | StartStage
	baseCommand := "adb shell am instrument -w -r com.example.swautoplay.test/androidx.test.runner.AndroidJUnitRunner"
	var params = []func(int, *gtk.Entry, *gtk.Entry, *gtk.Entry) (string, string, error){getAverageDungeonTime, getRunCount, getRefill, getDifficulty, getStartStage, getRunPosition, getDungeonName}
	index := dungeonsTabs.GetCurrentPage()
	log.Print(index)
	dungeon := dungeons[index]
	for i, fun := range params {
		if i >= len(dungeons) {
			name, value, err := fun(index, adt[index], runCount[index], startStage[index])
			if err == nil {
				baseCommand += " -e " + name + " " + value
			} else {
				log.Print("The parameter '" + name + "' is not filled")
			}
		} else {
			if dungeon.ConcernedParam[i] {
				name, value, err := fun(index, adt[index], runCount[index], startStage[index])
				if err == nil {
					baseCommand += " -e " + name + " " + value
				} else {
					log.Print("The parameter '" + name + "' is not filled")
				}
			}
		}
	}
	cmd := exec.Command(baseCommand)
	out, err := cmd.Output()
	if err != nil {
		log.Print("Error with adm command " + string(out))
	}
}

func getDungeonName(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return "DungeonName", dungeons[index].Name, nil
}

func getEntryText(entry *gtk.Entry, name string) (string, string, error) {
	value, err := entry.GetText()
	if value == "" {
		return name, "", fmt.Errorf("Empty string")
	}
	if err != nil {
		return name, "", err
	}
	return name, value, nil
}

func getAverageDungeonTime(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return getEntryText(adt, "AverageDungeonTime")
}

func getRunCount(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return getEntryText(runCount, "RunCount")
}

func getStartStage(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return getEntryText(startStage, "StartStage")
}

func getBoolParams(name string, params []*BoolProperty) (string, string, error) {
	for _, param := range params {
		if param.Value {
			return name, param.StringValue, nil
		}
	}
	return "", "", fmt.Errorf("bool param error")
}

func getDifficulty(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return getBoolParams("Difficulty", difficultyProps)
}

func getRefill(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return getBoolParams("Refill", refillProps)
}

func getRunPosition(index int, adt *gtk.Entry, runCount *gtk.Entry, startStage *gtk.Entry) (string, string, error) {
	return getBoolParams("RunPosition", props)
}

func (dungeon *Dungeon) createDungeonContent(count int, adt []*gtk.Entry, runCount []*gtk.Entry, startStage []*gtk.Entry) (*gtk.Grid, error) {
	contentGrid, err := gtk.GridNew()
	contentGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		log.Fatal("Unable to create dungeonsTabsChild:", err)
	}
	dungeonTitle, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create dungeonTitle:", err)
	}
	dungeonTitle.SetMarkup("<span size=\"large\" face=\"serif\"><b>Giant run parameters</b></span>")
	dungeonTitle.SetMarginBottom(15)
	dungeonTitle.SetMarginTop(10)
	dungeonTitle.SetHExpand(true)
	contentGrid.Add(dungeonTitle)

	if dungeon.ConcernedParam[0] {
		adt[count], _ = gtk.EntryNew()
		adtGrid, err := createGridEntry("Average dungeon time (in seconds) : ", 3, adt[count])
		if err != nil {
			log.Fatal("Unable to create adtGrid:", err)
		}
		contentGrid.Add(adtGrid)
	}

	if dungeon.ConcernedParam[1] {
		runCount[count], _ = gtk.EntryNew()
		runCountGrid, err := createGridEntry("Run count : ", 2, runCount[count])
		if err != nil {
			log.Fatal("Unable to create runCountGrid:", err)
		}
		contentGrid.Add(runCountGrid)
	}

	if dungeon.ConcernedParam[2] {
		refillGrid, err := createGridBoolBox("Refill energy from : ", refillProps)
		if err != nil {
			log.Fatal("Unable to create refillGrid:", err)
		}
		refillGrid.SetMarginTop(10)
		contentGrid.Add(refillGrid)
	}

	if dungeon.ConcernedParam[3] {
		difficultyGrid, err := createGridBoolBox(dungeon.Name+" difficulty : ", difficultyProps)
		if err != nil {
			log.Fatal("Unable to create difficultyGrid:", err)
		}
		difficultyGrid.SetMarginTop(10)
		contentGrid.Add(difficultyGrid)
	}

	if dungeon.ConcernedParam[4] {
		startStage[count], _ = gtk.EntryNew()
		startStageGrid, err := createGridEntry("Start dungeon to stage n° : ", 3, startStage[count])
		if err != nil {
			log.Fatal("Unable to create startStageGrid:", err)
		}
		contentGrid.Add(startStageGrid)
	}

	return contentGrid, nil
}

func createGridEntry(labelValue string, maxWidthChar int, entry *gtk.Entry) (*gtk.Grid, error) {
	entryGrid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	entryGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	entryGrid.SetMarginTop(10)
	entryLabel, err := createSubTitleLabel(labelValue)
	if err != nil {
		return nil, err
	}
	entryGrid.Add(entryLabel)

	entry.SetInputPurpose(gtk.INPUT_PURPOSE_NUMBER)
	entry.SetMaxWidthChars(maxWidthChar)
	entry.SetWidthChars(maxWidthChar)
	entryGrid.Add(entry)
	return entryGrid, nil
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
			log.Print(p.toString())
		})
		box.PackStart(radio[index], true, true, 0)
	}
	runPosGrid.Add(box)
	return runPosGrid, err
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
