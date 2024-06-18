package main

import (
	"log"

	"github.com/mickael-carl/git-tui/pkg/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}
