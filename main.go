package main

import (
	"dungeons-and-diagrams/editor"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {

	// Set log output to a file instead of stdout

	f, err := os.Create("run.log")
	if err != nil {
		log.Fatal("could not create log file", "error", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	p := tea.NewProgram(editor.New())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
