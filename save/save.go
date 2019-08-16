package save;

import (
  "os"
  "log"
  "io/ioutil"
  "io"
)

func ReadSaveFile(fileName string) string {
  path := "data/" + fileName
  if _, err := os.Stat(path); os.IsNotExist(err) {
    source, err := os.Open("data/savePattern/" + fileName)
    if err != nil {
      log.Fatal("Save src open error")
    }
    defer source.Close()

    destination, err := os.Create(path)
    if err != nil {
      log.Fatal("Save dest creation error")
    }
    defer destination.Close()
    _, err = io.Copy(destination, source)
    if err != nil {
      log.Fatal("Save copy error")
    }
  }
  content, err := ioutil.ReadFile(path)
  if err != nil {
    log.Fatal("Can't read save file")
  }
  return string(content)
}

func WriteSave(content string, fileName string) {
  err := ioutil.WriteFile("data/" + fileName, []byte(content), 0644)
	if err != nil {
		log.Fatal("Can't write in lastParams file")
	}
}