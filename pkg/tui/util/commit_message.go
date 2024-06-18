package util

import (
	"os"
	"os/exec"

	"github.com/rivo/tview"
)

func ReadCommitMessage(app *tview.Application, filename string) (string, error) {
	var err error
	app.Suspend(func() {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		cmd := exec.Command(editor, filename)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	})

	if err != nil {
		return "", err
	}

	message, err := os.ReadFile(filename)
	return string(message), err
}
