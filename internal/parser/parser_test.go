package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	file, err := os.Open("../../testdata/unknown-status.txt")
	//file, err := os.Open("../../testdata/panic.txt")
	//file, err := os.Open("../../testdata/few-errors.txt")
	if err != nil {
		t.Fatal(err)
	}

	events := ParseJsonsToGoTestEvents(file)

	assert.NotEmpty(t, events)

	events = TrimGoTestEvents(events)

	assert.NotEmpty(t, events)

	results := ExtractResults(events)

	assert.NotEmpty(t, results)

	outputFolder := "allure-results"

	CreateOutputFolder(outputFolder)

	PrintResults(outputFolder, results)

	assert.NotEmpty(t, results)
}
