package utils

import (
	"github.com/atotto/clipboard"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

// xsel required
func CopyToClipboard(value string) bool {
	err := clipboard.WriteAll(value)
	return err != nil
}