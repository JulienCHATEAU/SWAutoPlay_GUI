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
	cmd := exec.Command("adb", args...)
  joinedArgs := strings.Join(args, " ")
  fmt.Printf("|adb " + joinedArgs + "|\n")
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