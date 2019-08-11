package adb;

import (
  "os"
	"bytes"
	"io"
  "strings"
  "os/exec"
  "fmt"
)

func ExecAdbCommand(args ...string) (string) {
  buf := new(bytes.Buffer)
  joinedArgs := strings.Join(args, " ")
	cmd := exec.Command("adb", joinedArgs)
  fmt.Printf("adb " + joinedArgs + "\n")
	cmd.Stdout = io.MultiWriter(os.Stdout, buf)
  var stderr bytes.Buffer
  cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error : %s\nStderr : %s\n", err, stderr.String())
    return buf.String()
	}
  return buf.String()
}