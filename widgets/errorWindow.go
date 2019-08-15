package widgets

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func WaitForErrors(errors chan error, winChan chan *gtk.Window) {
	for {
		err := <-errors
		errorWindow := CreateErrorWindow(err)
		winChan <- errorWindow
	}
}

func CreateErrorWindow(err error) *gtk.Window {
	win, err2 := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err2 != nil {
		log.Panic("Can't create errorWindow (you became what you swore to handle !)")
	}
	win.SetDefaultSize(500, 60)
	win.SetTitle("Error")
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.Connect("destroy", gtk.MainQuit)

	erroGrid, err2 := gtk.GridNew()
	if err2 != nil {
		log.Fatal("Unable to create errorGrid:", err)
	}
	erroGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	erroGrid.SetMarginTop(10)
	erroGrid.SetMarginBottom(10)

	errText := fmt.Sprintf("%s", err)
	errorTitle, err2 := gtk.LabelNew("")
	if err2 != nil {
		log.Fatal("Unable to create errorTitle:", err)
	}
	errorTitle.SetMarkup("<span size=\"large\" face=\"serif\"><b>" + errText + "</b></span>")
	errorTitle.SetMarginBottom(15)
	errorTitle.SetMarginTop(10)
	errorTitle.SetHExpand(true)
	erroGrid.Add(errorTitle)

	win.Add(erroGrid)
	return win
}
