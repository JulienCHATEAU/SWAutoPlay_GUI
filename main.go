package main

import (
	"SWAutoPlay_GUI/adb"
	"SWAutoPlay_GUI/widgets"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"github.com/gotk3/gotk3/gtk"
	goadb "github.com/zach-klippenstein/goadb"
)

type AppWidgets struct {
	Adts          []*gtk.Entry
	RunCounts     []*gtk.Entry
	StartStages   []*gtk.Entry
	ScenarioNames []*gtk.ComboBoxText
}

type BoolProperty struct {
	Name        string
	Value       bool
	StringValue string
}

func (param *BoolProperty) toString() string {
	return fmt.Sprintf("Name : %s - Value : %t", param.Name, param.Value)
}

type Dungeon struct {
	Name            string
	ConcernedParam  [6]bool //AverageDungeonTime | RunCount | Refill | Difficulty | StartStage | HoH
	Adt             string
	RunCount        string
	StartStage      string
	ScenarioDungeon int
	BoolProps       [][]*BoolProperty // Refill | Difficulty | HoH
}

var scenarioDungeons = []string{"Garen", "Siz", "Kabir", "Ragon", "Telain", "Hydeni", "Tamor", "Vrofagus", "Faimon", "Aiden", "Ferun", "Runar", "Charuka"}
var startTestPosProps = []*BoolProperty{
	&BoolProperty{"Phone home page", true, "Home"},
	&BoolProperty{"Island", false, "Island"},
	&BoolProperty{"ToA stages page", false, "ToA"},
}

// var dungeons = []Dungeon{
// 	createDungeon("Giant", []bool{false, false, true}, "-1", "-1", "-1", 0, []int{3, 0, 1}),
// 	createDungeon("Drake", []bool{false, false, true}, "-1", "-1", "-1", 0, []int{3, 0, 1}),
// 	createDungeon("Necropolis", []bool{false, false, true}, "-1", "-1", "-1", 0, []int{3, 0, 1}),
// 	createDungeon("ToA", []bool{true, true, false}, "-1", "-1", "-1", 0, []int{3, 0, 1}),
// 	createDungeon("Scenario", []bool{true, true, false}, "-1", "-1", "-1", 5, []int{3, 0, 1}),
// }

var dungeons = createDungeonsFromSavedFile()

func (d *Dungeon) getActiveBoolProps() []string {
	boolPropsIndexes := make([]string, 3)
	for index, boolprop := range d.BoolProps {
		for i, prop := range boolprop {
			if prop.Value == true {
				boolPropsIndexes[index] = strconv.Itoa(i)
				break
			}
		}
	}
	return boolPropsIndexes
}

func (d *Dungeon) toSaveString(adt string, runCount string, startStage string, scenarioDungeon string) string {
	boolPropsIndexes := d.getActiveBoolProps()
	res := d.Name + "|" +
		strconv.FormatBool(d.ConcernedParam[3]) + "|" +
		strconv.FormatBool(d.ConcernedParam[4]) + "|" +
		strconv.FormatBool(d.ConcernedParam[5]) + "|" +
		adt + "|" +
		runCount + "|" +
		startStage + "|" +
		scenarioDungeon + "|" +
		boolPropsIndexes[0] + "|" +
		boolPropsIndexes[1] + "|" +
		boolPropsIndexes[2] + "\n"
	return res
}

func createDungeonsFromSavedFile() []Dungeon {
	dungeons := make([]Dungeon, 5)
	content, err := ioutil.ReadFile("data/lastParams")
	if err != nil {
		log.Fatal(err)
	}
	strContent := string(content)
	splitDungeon := strings.Split(strContent, "\n")
	for index, dungeon := range splitDungeon {
		if dungeon != "" {
			splitDungeonData := strings.Split(dungeon, "|")
			concernedParams := make([]bool, 3)
			for i := 1; i < 4; i++ {
				cp, err := strconv.ParseBool(splitDungeonData[i])
				if err != nil {
					if i == 1 {
						cp = true
					} else {
						cp = false
					}
				}
				concernedParams[i-1] = cp
			}
			scenarioDungeon, err := strconv.Atoi(splitDungeonData[7])
			if err != nil {
				scenarioDungeon = 0
			}
			radioSelectedIndex := make([]int, 3)
			for i := 8; i < len(splitDungeonData); i++ {
				rsi, err := strconv.Atoi(splitDungeonData[i])
				if err != nil {
					rsi = 0
				}
				radioSelectedIndex[i-8] = rsi
			}
			dungeons[index] = createDungeon(splitDungeonData[0],
				concernedParams,
				splitDungeonData[4],
				splitDungeonData[5],
				splitDungeonData[6],
				scenarioDungeon,
				radioSelectedIndex,
			)
		}
	}
	return dungeons
}

func createDungeon(name string, concernedParams []bool, adt string, runCount string, startStage string, scenarioDungeon int, radioSelectedIndex []int) Dungeon {
	hell := &BoolProperty{"Hell", false, "Hell"}
	dungeon := Dungeon{
		name,
		[6]bool{true, true, true, concernedParams[0], concernedParams[1], concernedParams[2]},
		adt,
		runCount,
		startStage,
		scenarioDungeon,
		[][]*BoolProperty{
			[]*BoolProperty{
				&BoolProperty{"Chest", false, "Chest"},
				&BoolProperty{"Social Point", false, "SocialPoint"},
				&BoolProperty{"Crystals", false, "Crystals"},
				&BoolProperty{"Don't refill", false, "Off"},
			},
			[]*BoolProperty{
				&BoolProperty{"Normal", false, "Normal"},
				&BoolProperty{"Hard", false, "Hard"},
			},
			[]*BoolProperty{
				&BoolProperty{"Yes", false, "true"},
				&BoolProperty{"No", false, "false"},
			},
		},
	}
	for index, prop := range dungeon.BoolProps {
		prop[radioSelectedIndex[index]].Value = true
	}
	if name != "ToA" {
		dungeon.BoolProps[1][2] = hell
	}
	return dungeon
}

func main() {

	gtk.Init(nil)

	dungeonsTabs, _ := gtk.NotebookNew()

	dunLength := len(dungeons)
	var appWidgets AppWidgets
	appWidgets.Adts = make([]*gtk.Entry, dunLength)
	appWidgets.RunCounts = make([]*gtk.Entry, dunLength)
	appWidgets.StartStages = make([]*gtk.Entry, dunLength)
	appWidgets.ScenarioNames = make([]*gtk.ComboBoxText, 1)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetDefaultSize(700, 550)
	win.SetTitle("SWAP")
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.Connect("destroy", func() {
		saveParams(appWidgets)
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
	if err != nil {
		log.Fatal("Unable to create notebook:", err)
	}
	dungeonsTabs.SetHExpand(true)
	dungeonsTabs.SetVExpand(true)

	for count, dungeon := range dungeons {
		contentGrid, _ := dungeon.createDungeonContent(count, appWidgets)

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
	runPosGrid, err := createGridBoolBox("Start this run from : ", startTestPosProps)
	if err != nil {
		log.Fatal("createGridBoolBox() failed :", err)
	}
	windowGrid.Add(runPosGrid)

	btnRun, err := gtk.ButtonNewWithLabel("Run")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btnStop, err := gtk.ButtonNewWithLabel("Exit")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btnRun.SetMarginTop(10)
	btnRun.SetMarginBottom(10)
	btnRun.SetMarginEnd(10)
	btnRun.SetMarginStart(10)
	btnStop.SetMarginBottom(10)
	btnStop.SetMarginEnd(10)
	btnStop.SetMarginStart(10)

	btnRun.Connect("clicked", func() {
		devices, _ := initDevices()
		runCommand := runCommand(dungeonsTabs, appWidgets)
		deviceWindow, err := widgets.CreateDeviceWindow(devices, runCommand, btnStop, btnRun)
		win.Connect("destroy", func() {
			deviceWindow.Close()
		})
		if err != nil {
			log.Print("Can't create device deviceWindow")
		}
		deviceWindow.ShowAll()
		gtk.Main()
	})
	btnStop.Connect("clicked", func() {
		saveParams(appWidgets)
		gtk.MainQuit()
		// gtk.Main()
	})
	btnStop.SetVisible(false)
	windowGrid.Add(btnRun)
	windowGrid.Add(btnStop)

	win.Add(windowGrid)
	win.ShowAll()

	gtk.Main()
}

func runCommand(dungeonsTabs *gtk.Notebook, appWidgets AppWidgets) []string { //AverageDungeonTime | RunCount | Refill | Difficulty | StartStage
	swautoplayPackage := "com.example.swautoplay.test/androidx.test.runner.AndroidJUnitRunner"
	args := []string{"instrument", "-w", "-r"}
	var params = []func(int, AppWidgets) (string, string, error){getAverageDungeonTime, getRunCount, getRefill, getDifficulty, getStartStage, getHoH, getRunPosition, getDungeonName}
	index := dungeonsTabs.GetCurrentPage()
	dungeon := dungeons[index]
	for i, fun := range params {
		if i >= len(dungeon.ConcernedParam) {
			name, value, err := fun(index, appWidgets)
			if err == nil {
				args = append(args, "-e", name, value)
			} else {
				log.Print("The parameter '" + name + "' is not filled")
			}
		} else {
			if dungeon.ConcernedParam[i] {
				name, value, err := fun(index, appWidgets)
				if err == nil {
					args = append(args, "-e", name, value)
				} else {
					log.Print("The parameter '" + name + "' is not filled")
				}
			}
		}
	}
	args = append(args, swautoplayPackage)
	return args
}

func initDevices() ([]widgets.Device, error) {
	out := adb.ExecAdbCommand("devices")
	outSplit := strings.Split(out, "\n")
	devices := make([]widgets.Device, len(outSplit)-3)
	for i := 1; outSplit[i] != "" && outSplit[i] != "\r"; i++ {
		deviceSplit := strings.Split(outSplit[i], "\t")
		devices[i-1].Serial = deviceSplit[0]
		if devices[i-1].IsWifi() {
			devices[i-1].Mode = "WiFi"
		} else {
			devices[i-1].Mode = "USB"
		}
		adb, err := goadb.New()
		if err != nil {
			return nil, err
		}
		devices[i-1].Manufacturer, err = adb.Device(goadb.DeviceWithSerial(devices[i-1].Serial)).RunCommand("getprop", "ro.product.manufacturer")
		devices[i-1].Model, err = adb.Device(goadb.DeviceWithSerial(devices[i-1].Serial)).RunCommand("getprop", "ro.product.model")
		devices[i-1].Manufacturer = strings.Trim(devices[i-1].Manufacturer, "\n")
		devices[i-1].Model = strings.Trim(devices[i-1].Model, "\n")
	}
	return devices, nil
}

func getDungeonName(index int, appWidgets AppWidgets) (string, string, error) {
	if dungeons[index].Name == "Scenario" {
		scenarioIndex := appWidgets.ScenarioNames[0].GetActive()
		return "DungeonName", scenarioDungeons[scenarioIndex], nil
	}
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

func getAverageDungeonTime(index int, appWidgets AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.Adts[index], "AverageDungeonTime")
}

func getRunCount(index int, appWidgets AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.RunCounts[index], "RunCount")
}

func getStartStage(index int, appWidgets AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.StartStages[index], "StartStage")
}

func getBoolParams(name string, params []*BoolProperty) (string, string, error) {
	for _, param := range params {
		if param.Value {
			return name, param.StringValue, nil
		}
	}
	return "", "", fmt.Errorf("bool param error")
}

func getRefill(index int, appWidgets AppWidgets) (string, string, error) {
	return getBoolParams("Refill", dungeons[index].BoolProps[0])
}

func getDifficulty(index int, appWidgets AppWidgets) (string, string, error) {
	return getBoolParams("Difficulty", dungeons[index].BoolProps[1])
}

func getHoH(index int, appWidgets AppWidgets) (string, string, error) {
	return getBoolParams("HoH", dungeons[index].BoolProps[2])
}

func getRunPosition(index int, appWidgets AppWidgets) (string, string, error) {
	return getBoolParams("StartTestPosition", startTestPosProps)
}

func saveParams(appWidgets AppWidgets) {
	content := ""
	for i, dungeon := range dungeons {
		startStage := ""
		adt, _ := appWidgets.Adts[i].GetText()
		runCount, _ := appWidgets.RunCounts[i].GetText()
		if dungeon.ConcernedParam[4] {
			startStage, _ = appWidgets.StartStages[i].GetText()
		}
		scenarioDungeon := strconv.Itoa(appWidgets.ScenarioNames[0].GetActive())
		content += dungeon.toSaveString(adt, runCount, startStage, scenarioDungeon)
	}
	err := ioutil.WriteFile("data/lastParams", []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (dungeon *Dungeon) createDungeonContent(count int, appWidgets AppWidgets) (*gtk.Grid, error) {
	contentGrid, err := gtk.GridNew()
	contentGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		log.Fatal("Unable to create dungeonsTabsChild:", err)
	}
	contentGrid.SetMarginTop(10)
	contentGrid.SetMarginBottom(10)

	dungeonTitle, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create dungeonTitle:", err)
	}
	dungeonTitle.SetMarkup("<span size=\"large\" face=\"serif\"><b>" + dungeon.Name + " run parameters</b></span>")
	dungeonTitle.SetMarginBottom(15)
	dungeonTitle.SetMarginTop(10)
	dungeonTitle.SetHExpand(true)
	contentGrid.Add(dungeonTitle)

	if dungeon.ConcernedParam[0] {
		appWidgets.Adts[count], _ = gtk.EntryNew()
		appWidgets.Adts[count].SetText(dungeon.Adt)
		adtGrid, err := createGridEntry("Average dungeon time (in seconds) : ", 3, appWidgets.Adts[count])
		if err != nil {
			log.Fatal("Unable to create adtGrid:", err)
		}
		contentGrid.Add(adtGrid)
	}

	if dungeon.ConcernedParam[1] {
		appWidgets.RunCounts[count], _ = gtk.EntryNew()
		appWidgets.RunCounts[count].SetText(dungeon.RunCount)
		runCountGrid, err := createGridEntry("Run count : ", 2, appWidgets.RunCounts[count])
		if err != nil {
			log.Fatal("Unable to create runCountGrid:", err)
		}
		contentGrid.Add(runCountGrid)
	}

	if dungeon.ConcernedParam[2] {
		refillGrid, err := createGridBoolBox("Refill energy from : ", dungeon.BoolProps[0])
		if err != nil {
			log.Fatal("Unable to create refillGrid:", err)
		}
		contentGrid.Add(refillGrid)
	}

	if dungeon.Name == "Scenario" {
		appWidgets.ScenarioNames[0], _ = gtk.ComboBoxTextNew()
		for _, name := range scenarioDungeons {
			appWidgets.ScenarioNames[0].AppendText(name)
		}
		appWidgets.ScenarioNames[0].SetActive(dungeon.ScenarioDungeon)

		boxGrid, err := gtk.GridNew()
		if err != nil {
			return nil, err
		}
		boxGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
		boxGrid.SetMarginTop(10)
		entryLabel, err := createSubTitleLabel("Scenario dungeon : ")
		if err != nil {
			return nil, err
		}
		boxGrid.Add(entryLabel)
		boxGrid.Add(appWidgets.ScenarioNames[0])
		contentGrid.Add(boxGrid)
	}

	if dungeon.ConcernedParam[3] {
		difficultyGrid, err := createGridBoolBox(dungeon.Name+" difficulty : ", dungeon.BoolProps[1])
		if err != nil {
			log.Fatal("Unable to create difficultyGrid:", err)
		}
		contentGrid.Add(difficultyGrid)
	}

	if dungeon.ConcernedParam[4] {
		appWidgets.StartStages[count], _ = gtk.EntryNew()
		appWidgets.StartStages[count].SetText(dungeon.StartStage)
		startStageGrid, err := createGridEntry("Start dungeon to stage n° : ", 3, appWidgets.StartStages[count])
		if err != nil {
			log.Fatal("Unable to create startStageGrid:", err)
		}
		contentGrid.Add(startStageGrid)
	}

	if dungeon.ConcernedParam[5] {
		hohGrid, err := createGridBoolBox("Is there any opened HoH : ", dungeon.BoolProps[2])
		if err != nil {
			log.Fatal("Unable to create hohGrid:", err)
		}
		contentGrid.Add(hohGrid)
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
