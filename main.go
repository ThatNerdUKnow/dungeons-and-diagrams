package main

import (
	"dungeons-and-diagrams/model"
	"fmt"
	"os"
	"runtime/pprof"

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

	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}

	defer cpuFile.Close()
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()
	p := tea.NewProgram(model.New())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
