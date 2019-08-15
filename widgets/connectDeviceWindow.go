package widgets;

import (
  "fmt"
	"github.com/gotk3/gotk3/gtk"
  "regexp"
  "SWAutoPlay_GUI/adb"
  "strings"
)

func CreateConnectDeviceWindow(win *gtk.Window, errChan chan error, winChan chan *gtk.Window) (*gtk.Window, error) {
  window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
    return nil, err
	}
	window.SetDefaultSize(430, 150)
	window.SetTitle("Device connection")
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.Connect("destroy", gtk.MainQuit)

	contentGrid, err := gtk.GridNew()
	if err != nil {
		return nil, fmt.Errorf("Unable to create contentGrid:", err)
	}
	contentGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)

  title, err := gtk.LabelNew("")
	if err != nil {
		return nil, fmt.Errorf("Unable to create errorTitle:", err)
	}
	title.SetMarkup("<span size=\"large\" face=\"serif\"><b>Connect a new device over WiFi</b></span>")
	title.SetMarginBottom(15)
	title.SetMarginTop(10)
	title.SetHExpand(true)
	contentGrid.Add(title)

  entry, _ := gtk.EntryNew()
  entry.SetText("192.168.1.20")
  ipAddressGrid, err := CreateGridEntry("IP address of the device to connect : ", 16, entry)
  ipAddressGrid.SetMarginStart(10)
  ipAddressGrid.SetMarginEnd(20)
  if err != nil {
    return nil, fmt.Errorf("Unable to create ipAddressGrid:", err)
  }
  contentGrid.Add(ipAddressGrid)

  connectButton, err := gtk.ButtonNewWithLabel("Connect")
	if err != nil {
		return nil, fmt.Errorf("Unable to create connectButton:", err)
	}
	connectButton.SetMarginTop(25)
	connectButton.SetMarginBottom(10)
	connectButton.SetMarginEnd(10)
	connectButton.SetMarginStart(10)
  connectButton.Connect("clicked", func() {
    ip, _ := entry.GetText()
    err := connectWithIpAddress(ip)
    if err != nil {
      HandleError(win, errChan, winChan, err)
    } else {
      window.Close()
    }
  })

  contentGrid.Add(connectButton)

	window.Add(contentGrid)
  return window, nil
}

func connectWithIpAddress(ip string) error {
  ipMatch := checkIpAddress(ip)
  if !ipMatch {
    return fmt.Errorf("That's not an IP address")
  } else {
    connectionAddress := ip+":5555"
    out := strings.TrimRight(adb.ExecAdbCommand("connect", connectionAddress), "\r\n")
    switch out {
    case "already connected to "+ connectionAddress:
      return fmt.Errorf("This device is already connected over WiFi")
      break
    case "connected to "+ connectionAddress:
      return nil
      break
    default:
      return fmt.Errorf("Can't connect to this device (port 5555 might not be opened)")
      break
    }
  }
  return nil
}

func checkIpAddress(ipAddress string) bool {
   ipAddress = strings.Trim(ipAddress, " ")

   re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
   if re.MatchString(ipAddress) {
           return true
   }
   return false
}