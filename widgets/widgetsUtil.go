package widgets;

import (
  "github.com/gotk3/gotk3/gtk"
  "log"
  "fmt"
  "strconv"
)
var ScenarioDungeons = []string{"Garen", "Siz", "Kabir", "Ragon", "Telain", "Hydeni", "Tamor", "Vrofagus", "Faimon", "Aiden", "Ferun", "Runar", "Charuka"}

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

func (d *Dungeon) ToSaveString(adt string, runCount string, startStage string, scenarioDungeon string) string {
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

func (dungeon *Dungeon) CreateDungeonContent(count int, appWidgets AppWidgets) (*gtk.Grid, error) {
	contentGrid, err := gtk.GridNew()
	contentGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		log.Fatal("Unable to Create dungeonsTabsChild:", err)
	}
	contentGrid.SetMarginTop(10)
	contentGrid.SetMarginBottom(10)

	dungeonTitle, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to Create dungeonTitle:", err)
	}
	dungeonTitle.SetMarkup("<span size=\"large\" face=\"serif\"><b>" + dungeon.Name + " run parameters</b></span>")
	dungeonTitle.SetMarginBottom(15)
	dungeonTitle.SetMarginTop(10)
	dungeonTitle.SetHExpand(true)
	contentGrid.Add(dungeonTitle)

	if dungeon.ConcernedParam[0] {
		appWidgets.Adts[count], _ = gtk.EntryNew()
		appWidgets.Adts[count].SetText(dungeon.Adt)
		adtGrid, err := CreateGridEntry("Average dungeon time (in seconds) : ", 3, appWidgets.Adts[count])
		if err != nil {
			log.Fatal("Unable to Create adtGrid:", err)
		}
		contentGrid.Add(adtGrid)
	}

	if dungeon.ConcernedParam[1] {
		appWidgets.RunCounts[count], _ = gtk.EntryNew()
		appWidgets.RunCounts[count].SetText(dungeon.RunCount)
		runCountGrid, err := CreateGridEntry("Run count : ", 2, appWidgets.RunCounts[count])
		if err != nil {
			log.Fatal("Unable to Create runCountGrid:", err)
		}
		contentGrid.Add(runCountGrid)
	}

	if dungeon.ConcernedParam[2] {
		refillGrid, err := CreateGridBoolBox("Refill energy from : ", dungeon.BoolProps[0])
		if err != nil {
			log.Fatal("Unable to Create refillGrid:", err)
		}
		contentGrid.Add(refillGrid)
	}

	if dungeon.Name == "Scenario" {
		appWidgets.ScenarioNames[0], _ = gtk.ComboBoxTextNew()
		for _, name := range ScenarioDungeons {
			appWidgets.ScenarioNames[0].AppendText(name)
		}
		appWidgets.ScenarioNames[0].SetActive(dungeon.ScenarioDungeon)

		boxGrid, err := gtk.GridNew()
		if err != nil {
			return nil, err
		}
		boxGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
		boxGrid.SetMarginTop(10)
		entryLabel, err := CreateSubTitleLabel("Scenario dungeon : ")
		if err != nil {
			return nil, err
		}
		boxGrid.Add(entryLabel)
		boxGrid.Add(appWidgets.ScenarioNames[0])
		contentGrid.Add(boxGrid)
	}

	if dungeon.ConcernedParam[3] {
		difficultyGrid, err := CreateGridBoolBox(dungeon.Name+" difficulty : ", dungeon.BoolProps[1])
		if err != nil {
			log.Fatal("Unable to Create difficultyGrid:", err)
		}
		contentGrid.Add(difficultyGrid)
	}

	if dungeon.ConcernedParam[4] {
		appWidgets.StartStages[count], _ = gtk.EntryNew()
		appWidgets.StartStages[count].SetText(dungeon.StartStage)
		startStageGrid, err := CreateGridEntry("Start dungeon to stage n° : ", 3, appWidgets.StartStages[count])
		if err != nil {
			log.Fatal("Unable to Create startStageGrid:", err)
		}
		contentGrid.Add(startStageGrid)
	}

	if dungeon.ConcernedParam[5] {
		hohGrid, err := CreateGridBoolBox("Is there any opened HoH : ", dungeon.BoolProps[2])
		if err != nil {
			log.Fatal("Unable to Create hohGrid:", err)
		}
		contentGrid.Add(hohGrid)
	}

	return contentGrid, nil
}





func CreateGridEntry(labelValue string, maxWidthChar int, entry *gtk.Entry) (*gtk.Grid, error) {
	entryGrid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	entryGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	entryGrid.SetMarginTop(10)
	entryLabel, err := CreateSubTitleLabel(labelValue)
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

func CreateGridBoolBox(labelValue string, props []*BoolProperty) (*gtk.Grid, error) {
	runPosGrid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	runPosGrid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	runPosLabel, err := CreateSubTitleLabel(labelValue)
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

func updateParam(param *BoolProperty, state bool) {
	param.Value = state
}

func CreateSubTitleLabel(name string) (*gtk.Label, error) {
	label, err := gtk.LabelNew(name)
	label.SetMarginBottom(5)
	label.SetMarginTop(5)
	label.SetMarginStart(10)
	label.SetMarginEnd(25)
	return label, err
}

func HandleError(mainWin *gtk.Window, errorsChan chan error, winChan chan *gtk.Window, err error) {
	errorWindow := CreateErrorWindow(err)
	mainWin.Connect("destroy", func() {
		errorWindow.Close()
		gtk.MainQuit()
	})
	errorWindow.ShowAll()
	gtk.Main()
}
