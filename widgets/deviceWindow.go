package widgets

import (
	"fmt"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	goadb "github.com/zach-klippenstein/goadb"
)

// var devices = []Device{
//   Device{"Sony", "F1662", "192.168.1.20", "WiFi"},
//   Device{"Huawei", "X8976", "FEF698ST65", "USB"},
// }

var selectedDeviceIndex = -1

type Device struct {
	Manufacturer string
	Model        string
	Serial       string
	Mode         string
}

func (device *Device) IsWifi() bool {
	return strings.Contains(device.Serial, ".")
}

func (device *Device) ToLabel() string {
	return device.Manufacturer + " " + device.Model + "  (" + device.Serial + ") - " + device.Mode
}

func (device *Device) IsItMyLabel(label string) bool {
	return device.ToLabel() == label
}

func CreateDeviceWindow(devices []Device, runCommand []string, btnStop *gtk.Button, btnRun *gtk.Button) (*gtk.Window, error) {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	win.SetDefaultSize(400, 280)
	win.SetTitle("Device choice")
	win.SetPosition(gtk.WIN_POS_CENTER)

	contentGrid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	contentGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	contentTitle, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create contentTitle:", err)
	}
	contentTitle.SetMarkup("<span size=\"large\" face=\"serif\"><b>Chose your device</b></span>")
	contentTitle.SetMarginBottom(15)
	contentTitle.SetMarginTop(10)
	contentTitle.SetHExpand(true)
	contentGrid.Add(contentTitle)

	treeView, _ := gtk.TreeViewNew()
	treeView.SetHExpand(true)
	treeView.SetVExpand(true)
	listStore, _ := gtk.ListStoreNew(glib.TYPE_STRING)

	// Window properties
	win.SetTitle("Device choice")
	win.Connect("destroy", gtk.MainQuit)

	// treeView properties
	renderer, _ := gtk.CellRendererTextNew()
	column, _ := gtk.TreeViewColumnNewWithAttribute("Device name", renderer, "text", 0)
	treeView.AppendColumn(column)
	treeView.SetModel(listStore)

	for _, device := range devices {
		listStore.SetValue(listStore.Append(), 0, device.ToLabel())
	}

	// treeView selection properties
	sel, _ := treeView.GetSelection()
	sel.Connect("changed", func() {
		deviceChanged(sel, listStore, devices)
	})
	contentGrid.Add(treeView)

	//run button
	btn, err := gtk.ButtonNewWithLabel("Run")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.SetMarginTop(10)
	btn.SetMarginBottom(10)
	btn.SetMarginEnd(10)
	btn.SetMarginStart(10)
	btn.Connect("clicked", func() {
		if len(devices) > 0 {
			win.Destroy()
			go func() {
				run(devices[selectedDeviceIndex], runCommand)
				btnRun.SetVisible(true)
				btnStop.SetLabel("Exit")
			}()
			btnRun.SetVisible(false)
			btnStop.SetLabel("Stop current run")
		}
	})
	contentGrid.Add(btn)

	win.Add(contentGrid)
	return win, nil
}

func run(device Device, runCommand []string) {
	adb, err := goadb.New()
	if err != nil {
		fmt.Printf("Error creating adb %s", err)
	}
	fmt.Printf("%q\n", runCommand)
	out, err := adb.Device(goadb.DeviceWithSerial(device.Serial)).RunCommand("am", runCommand...)
	if err != nil {
		fmt.Printf("Error with adb command %s", err)
	}
	fmt.Printf(out)
}

func deviceChanged(s *gtk.TreeSelection, listStore *gtk.ListStore, devices []Device) error {
	rows := s.GetSelectedRows(listStore)
	var value *glib.Value
	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		iter, _ := listStore.GetIter(path)
		value, _ = listStore.GetValue(iter, 0)
		deviceLabel, _ := value.GetString()
		for i, device := range devices {
			if device.IsItMyLabel(deviceLabel) {
				selectedDeviceIndex = i
				break
			}
		}
	}
	fmt.Printf("selectedDeviceIndex : %d", selectedDeviceIndex)
	return nil
}
