package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyubin/gotest2allure/internal/parser"
)

var (
	fileFlag         string
	outputFlag       string
	issuePatternFrag string
)

const (
	fileFlagUsage       = "tests results file"
	outputFlagUsage     = "output results folder"
	outputFlagDefault   = "allure-results"
	issuePatternUsage   = "patter issue link for bug tracker"
	issuePatternDefault = "%s"
)

func main() {

	flag.StringVar(&fileFlag, "file", "", fileFlagUsage)
	flag.StringVar(&fileFlag, "f", "", fileFlagUsage)
	flag.StringVar(&outputFlag, "output", outputFlagDefault, outputFlagUsage)
	flag.StringVar(&outputFlag, "o", outputFlagDefault, outputFlagUsage)
	flag.StringVar(&issuePatternFrag, "issuePattern", issuePatternDefault, issuePatternUsage)

	flag.Parse()

	if fileFlag == "" {
		panic("no file")
	}

	f, err := os.Open(filepath.Clean(fileFlag))
	if fileFlag == "" {
		panic(err)
	}

	if !strings.Contains(issuePatternFrag, "%s") {
		fmt.Println("warning: issuePattern frag shoud have %s symbol")
		issuePatternFrag = "%s"
	}

	events := parser.TrimGoTestEvents(parser.ParseJsonsToGoTestEvents(f))
	results := parser.ExtractResults(events, issuePatternFrag)

	parser.CreateOutputFolder(outputFlag)
	parser.PrintResults(outputFlag, results)
}
