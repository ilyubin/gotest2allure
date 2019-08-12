package main

import (
	"flag"
	"github.com/ilyubin/gotest2allure/internal/parser"
	"os"
	"path/filepath"
)

var (
	fileFlag   string
	outputFlag string
)

const (
	fileFlagUsage     = "tests results file"
	outputFlagUsage   = "output results folder"
	defaultOutputFlag = "allure-results"
)

func main() {

	flag.StringVar(&fileFlag, "file", "", fileFlagUsage)
	flag.StringVar(&fileFlag, "f", "", fileFlagUsage)
	flag.StringVar(&outputFlag, "output", defaultOutputFlag, outputFlagUsage)
	flag.StringVar(&outputFlag, "o", defaultOutputFlag, outputFlagUsage)

	flag.Parse()

	if fileFlag == "" {
		panic("no file")
	}

	f, err := os.Open(filepath.Clean(fileFlag))
	if fileFlag == "" {
		panic(err)
	}

	events := parser.TrimGoTestEvents(parser.ParseJsonsToGoTestEvents(f))
	containers := parser.ExtractContainers(events)
	results := parser.ExtractResults(events, containers)

	parser.CreateOutputFolder(outputFlag)
	parser.PrintContainers(outputFlag, containers)
	parser.PrintResults(outputFlag, results)
}
