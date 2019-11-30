package main

import (
	"SWAutoPlay_GUI/adb"
	"SWAutoPlay_GUI/save"
	wid "SWAutoPlay_GUI/widgets"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

const DUNGEON_COUNT = 8

var startTestPosProps = []*wid.BoolProperty{
	&wid.BoolProperty{"Phone home page", true, "Home"},
	&wid.BoolProperty{"Island", false, "Island"},
	&wid.BoolProperty{"ToA stages page", false, "ToA"},
}

var dungeons = createDungeonsFromSavedFile()

func main() {

	gtk.Init(nil)

	errorsChan := make(chan error)
	winChan := make(chan *gtk.Window)
	go wid.WaitForErrors(errorsChan, winChan)
	dungeonsTabs, _ := gtk.NotebookNew()

	dunLength := len(dungeons)
	var appWidgets wid.AppWidgets
	appWidgets.Adts = make([]*gtk.Entry, dunLength)
	appWidgets.RunCounts = make([]*gtk.Entry, dunLength)
	appWidgets.StartStages = make([]*gtk.Entry, dunLength)
	appWidgets.Level = make([]*gtk.Entry, dunLength)
	appWidgets.ScenarioNames = make([]*gtk.ComboBoxText, 1)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetDefaultSize(700, 550)
	win.SetTitle("SWAP")
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.Connect("destroy", func() {
		SaveParams(appWidgets)
		win.Destroy()
		gtk.MainQuit()
		disconnectDevice(wid.CurrentRunSerial)
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
		contentGrid, _ := dungeon.CreateDungeonContent(count, appWidgets)

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
	runPosGrid, err := wid.CreateGridBoolBox("Start this run from : ", startTestPosProps)
	runPosGrid.SetMarginTop(10)
	runPosGrid.SetMarginStart(10)
	if err != nil {
		log.Fatal("createGridBoolBox() failed :", err)
	}
	windowGrid.Add(runPosGrid)

	buttonGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create buttonGrid:", err)
	}
	buttonGrid.SetHExpand(true)
	buttonGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	buttonGrid.SetMarginTop(10)
	buttonGrid.SetMarginBottom(10)
	buttonGrid.SetMarginEnd(10)
	buttonGrid.SetMarginStart(10)

	btnRun, err := gtk.ButtonNewWithLabel("Run this dungeon")
	if err != nil {
		log.Fatal("Unable to create btnRun:", err)
	}
	btnRun.SetHExpand(true)
	btnRun.SetMarginEnd(10)
	btnStop, err := gtk.ButtonNewWithLabel("Exit")
	if err != nil {
		log.Fatal("Unable to create btnStop:", err)
	}
	btnStop.SetHExpand(true)
	btnConnect, err := gtk.ButtonNewWithLabel("Connect new devices")
	if err != nil {
		log.Fatal("Unable to create btnConnect:", err)
	}
	btnConnect.SetHExpand(true)
	btnConnect.SetMarginEnd(10)

	btnConnect.Connect("clicked", func() {
		connectDeviceWindow, err := wid.CreateConnectDeviceWindow(win, errorsChan, winChan)
		if err != nil {
			fmt.Printf("Error : ", err)
		}
		win.Connect("destroy", func() {
			connectDeviceWindow.Close()
			gtk.MainQuit()
		})
		connectDeviceWindow.ShowAll()
		gtk.Main()
	})

	btnRun.Connect("clicked", func() {
		devices, _ := initDevices()
		runCommand, err := createRunCommand(dungeonsTabs, appWidgets)
		if err != nil {
			wid.HandleError(win, errorsChan, winChan, err)
		} else {
			devicesWindow, err := wid.CreateDeviceWindow(win, devices, runCommand, btnStop, btnRun)
			if err != nil {
				log.Print("Can't create device devicesWindow")
			}
			win.Connect("destroy", func() {
				devicesWindow.Close()
				gtk.MainQuit()
			})
			devicesWindow.ShowAll()
			gtk.Main()
		}
	})
	btnStop.Connect("clicked", func() {
		if value, _ := btnStop.GetLabel(); value == "Exit" {
			SaveParams(appWidgets)
			gtk.MainQuit()
		} else {
			stopRun()
		}
		// gtk.Main()
	})
	btnStop.SetVisible(false)
	buttonGrid.Add(btnRun)
	buttonGrid.Add(btnConnect)
	buttonGrid.Add(btnStop)
	windowGrid.Add(buttonGrid)

	win.Add(windowGrid)
	win.ShowAll()

	gtk.Main()
}

func stopRun() {
	disconnectDevice(wid.CurrentRunSerial)
	time.Sleep(1 * time.Second)
	adb.ExecAdbCommand("connect", wid.CurrentRunSerial)
}

func disconnectDevice(serial string) {
	adb.ExecAdbCommand("disconnect", serial)
}

func createRunCommand(dungeonsTabs *gtk.Notebook, appWidgets wid.AppWidgets) ([]string, error) { //AverageDungeonTime | RunCount | Refill | Difficulty | StartStage
	swautoplayPackage := "com.example.swautoplay.test/androidx.test.runner.AndroidJUnitRunner"
	args := []string{"instrument", "-w", "-r"}
	var params = []func(int, wid.AppWidgets) (string, string, error){getAverageDungeonTime, getRunCount, getRefill, getDifficulty, getStartStage, getHoH, getLevel, getRunPosition, getDungeonName}
	index := dungeonsTabs.GetCurrentPage()
	dungeon := dungeons[index]
	for i, fun := range params {
		if i >= len(dungeon.ConcernedParam) {
			name, value, err := fun(index, appWidgets)
			if err == nil {
				args = append(args, "-e", name, value)
			} else {
				return nil, err
			}
		} else {
			if dungeon.ConcernedParam[i] {
				name, value, err := fun(index, appWidgets)
				if err == nil {
					args = append(args, "-e", name, value)
				} else {
					return nil, err
				}
			}
		}
	}
	args = append(args, swautoplayPackage)
	return args, nil
}

func SaveParams(appWidgets wid.AppWidgets) {
	content := ""
	for i, dungeon := range dungeons {
		startStage := ""
		adt, _ := appWidgets.Adts[i].GetText()
		level, _ := appWidgets.Level[i].GetText()
		runCount, _ := appWidgets.RunCounts[i].GetText()
		if dungeon.ConcernedParam[4] {
			startStage, _ = appWidgets.StartStages[i].GetText()
		}
		scenarioDungeon := strconv.Itoa(appWidgets.ScenarioNames[0].GetActive())
		content += dungeon.ToSaveString(adt, runCount, startStage, scenarioDungeon, level)
	}
	save.WriteSave(content, "lastParams")
}

func initDevices() ([]wid.Device, error) {
	out := adb.ExecAdbCommand("devices")
	outSplit := strings.Split(out, "\n")
	devices := make([]wid.Device, len(outSplit)-3)
	if strings.HasPrefix(outSplit[1], "* daemon") {
		return make([]wid.Device, 0), nil
	}
	for i := 1; outSplit[i] != "" && outSplit[i] != "\r"; i++ {
		deviceSplit := strings.Split(outSplit[i], "\t")
		devices[i-1].Serial = deviceSplit[0]
		if devices[i-1].IsWifi() {
			devices[i-1].Mode = "WiFi"
		} else {
			devices[i-1].Mode = "USB"
		}
		devices[i-1].Manufacturer = strings.TrimRight(adb.ExecAdbCommand("-s", devices[i-1].Serial, "shell", "getprop", "ro.product.manufacturer"), "\r\n")
		devices[i-1].Model = strings.TrimRight(adb.ExecAdbCommand("-s", devices[i-1].Serial, "shell", "getprop", "ro.product.model"), "\r\n")
	}
	return devices, nil
}

func getLevel(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.Level[index], "Level")
}

func getDungeonName(index int, appWidgets wid.AppWidgets) (string, string, error) {
	if dungeons[index].Name == "Scenario" {
		scenarioIndex := appWidgets.ScenarioNames[0].GetActive()
		return "DungeonName", wid.ScenarioDungeons[scenarioIndex], nil
	}
	return "DungeonName", dungeons[index].Name, nil
}

func getEntryText(entry *gtk.Entry, name string) (string, string, error) {
	value, err := entry.GetText()
	if value == "" {
		return name, "", fmt.Errorf(name + " entry is empty")
	}
	if err != nil {
		return name, "", err
	}
	if _, err := strconv.Atoi(value); err != nil {
		return name, "", fmt.Errorf(name + " entry must be a number")
	}
	return name, value, nil
}

func getAverageDungeonTime(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.Adts[index], "AverageDungeonTime")
}

func getRunCount(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.RunCounts[index], "RunCount")
}

func getStartStage(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getEntryText(appWidgets.StartStages[index], "StartStage")
}

func getBoolParams(name string, params []*wid.BoolProperty) (string, string, error) {
	for _, param := range params {
		if param.Value {
			return name, param.StringValue, nil
		}
	}
	return "", "", fmt.Errorf("bool param error")
}

func getRefill(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getBoolParams("Refill", dungeons[index].BoolProps[0])
}

func getDifficulty(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getBoolParams("Difficulty", dungeons[index].BoolProps[1])
}

func getHoH(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getBoolParams("HoH", dungeons[index].BoolProps[2])
}

func getRunPosition(index int, appWidgets wid.AppWidgets) (string, string, error) {
	return getBoolParams("StartTestPosition", startTestPosProps)
}

func createDungeonsFromSavedFile() []wid.Dungeon {
	dungeons := make([]wid.Dungeon, DUNGEON_COUNT)
	strContent := save.ReadSaveFile("lastParams")
	splitDungeon := strings.Split(strContent, "\n")
	dcbpc := wid.DUNGEON_CONCERNED_BOOL_PARAM_COUNT
	for index, dungeon := range splitDungeon {
		if dungeon != "" {
			splitDungeonData := strings.Split(dungeon, "|")
			concernedParams := make([]bool, dcbpc)
			for i := 1; i < dcbpc+1; i++ {
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
			scenarioDungeon, err := strconv.Atoi(splitDungeonData[dcbpc+5])
			if err != nil {
				scenarioDungeon = 0
			}
			radioSelectedIndex := make([]int, 3)
			startt := dcbpc + 6
			for i := startt; i < len(splitDungeonData); i++ {
				rsi, err := strconv.Atoi(splitDungeonData[i])
				if err != nil {
					rsi = 0
				}
				radioSelectedIndex[i-startt] = rsi
			}
			dungeons[index] = createDungeon(splitDungeonData[0],
				concernedParams,
				splitDungeonData[5],
				splitDungeonData[6],
				splitDungeonData[7],
				splitDungeonData[8],
				scenarioDungeon,
				radioSelectedIndex,
			)
		}
	}
	return dungeons
}

func createDungeon(name string, concernedParams []bool, adt string, runCount string, startStage string, level string, scenarioDungeon int, radioSelectedIndex []int) wid.Dungeon {
	hell := &wid.BoolProperty{"Hell", false, "Hell"}
	difficulty_props := []*wid.BoolProperty{
		&wid.BoolProperty{"Normal", false, "Normal"},
		&wid.BoolProperty{"Hard", false, "Hard"},
	}
	if (wid.IsRiftDungeon(name)) {
		zone := "Forest"
		if (name == "Ellunia") {
			zone = "Sanctuary"
		} else if (name == "Lumel") {
			zone = "Cliff"
		}
		difficulty_props = []*wid.BoolProperty{
			&wid.BoolProperty{"Vestige" , false, "Vestige"},
			&wid.BoolProperty{zone, false, "Rune"},
		}
	}
	dungeon := wid.Dungeon{
		name,
		[wid.DUNGEON_CONCERNED_PARAM_COUNT]bool{true, true, true, concernedParams[0], concernedParams[1], concernedParams[2], concernedParams[3]},
		adt,
		runCount,
		startStage,
		level,
		scenarioDungeon,
		[][]*wid.BoolProperty{
			[]*wid.BoolProperty{
				&wid.BoolProperty{"Chest", false, "Chest"},
				&wid.BoolProperty{"Social Point", false, "SocialPoint"},
				&wid.BoolProperty{"Crystals", false, "Crystals"},
				&wid.BoolProperty{"Don't refill", false, "Off"},
			},
			difficulty_props,
			[]*wid.BoolProperty{
				&wid.BoolProperty{"Yes", false, "true"},
				&wid.BoolProperty{"No", false, "false"},
			},
		},
	}
	if name != "ToA" && !wid.IsRiftDungeon(name) {
		dungeon.BoolProps[1] = append(dungeon.BoolProps[1], hell)
	}
	for index, prop := range dungeon.BoolProps {
		prop[radioSelectedIndex[index]].Value = true
	}
	return dungeon
}
