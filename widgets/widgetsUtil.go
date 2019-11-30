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
	Level 		  []*gtk.Entry
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

const DUNGEON_CONCERNED_PARAM_COUNT = 7
const DUNGEON_CONCERNED_BOOL_PARAM_COUNT = 4 // Difficulty | StartStage | HoH | Level

type Dungeon struct {
	Name            string
	ConcernedParam  [DUNGEON_CONCERNED_PARAM_COUNT]bool //AverageDungeonTime | RunCount | Refill | Difficulty | StartStage | HoH | Level
	Adt             string
	RunCount        string
	StartStage      string
	Level 			string
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

func (d *Dungeon) ToSaveString(adt string, runCount string, startStage string, scenarioDungeon string, level string) string {
	boolPropsIndexes := d.getActiveBoolProps()
	res := d.Name + "|" +
		strconv.FormatBool(d.ConcernedParam[3]) + "|" +
		strconv.FormatBool(d.ConcernedParam[4]) + "|" +
		strconv.FormatBool(d.ConcernedParam[5]) + "|" +
		strconv.FormatBool(d.ConcernedParam[6]) + "|" +
		adt + "|" +
		runCount + "|" +
		startStage + "|" +
		level + "|" +
		scenarioDungeon + "|" +
		boolPropsIndexes[0] + "|" +
		boolPropsIndexes[1] + "|" +
		boolPropsIndexes[2] + "\n"
	return res
}

func IsRiftDungeon(dungeonName string) bool {
	return dungeonName == "Karzhan" || dungeonName == "Ellunia" || dungeonName == "Lumel"
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
		boxGrid.SetMarginBottom(10)
		entryLabel, err := CreateSubTitleLabel("Scenario dungeon : ")
		if err != nil {
			return nil, err
		}
		boxGrid.Add(entryLabel)
		boxGrid.Add(appWidgets.ScenarioNames[0])
		contentGrid.Add(boxGrid)
	}

	if dungeon.ConcernedParam[3] {
		difficultyStr := " difficulty : "
		if (IsRiftDungeon(dungeon.Name)) {
			difficultyStr = " dungeon : "
		}
		difficultyGrid, err := CreateGridBoolBox(dungeon.Name + difficultyStr, dungeon.BoolProps[1])
		if err != nil {
			log.Fatal("Unable to Create difficultyGrid:", err)
		}
		contentGrid.Add(difficultyGrid)
	}

	if dungeon.ConcernedParam[4] {
		appWidgets.StartStages[count], _ = gtk.EntryNew()
		appWidgets.StartStages[count].SetText(dungeon.StartStage)
		label := "Start dungeon to stage n° : "
		if (IsRiftDungeon(dungeon.Name)) {
			firstMonster := "Inugami"
			secondMonster := "Bear"
			if (dungeon.Name == "Ellunia") {
				firstMonster = "Fairy"
				secondMonster = "Pixie"
			} else if (dungeon.Name == "Lumel") {
				firstMonster = "Werewolf"
				secondMonster = "Martial cat"
			}
			label = "Monster (0 for "+firstMonster+" and 1 for "+secondMonster+") :"
		}
		startStageGrid, err := CreateGridEntry(label, 3, appWidgets.StartStages[count])
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

	if dungeon.ConcernedParam[6] {
		appWidgets.Level[count], _ = gtk.EntryNew()
		appWidgets.Level[count].SetText(dungeon.Level)
		levelGrid, err := CreateGridEntry("Dungeon level : ", 3, appWidgets.Level[count])
		if err != nil {
			log.Fatal("Unable to Create levelGrid:", err)
		}
		contentGrid.Add(levelGrid)
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
  runPosGrid.SetMarginTop(10)
	runPosGrid.SetMarginBottom(10)
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
