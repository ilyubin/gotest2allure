package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	file, err := os.Open("../../testdata/unknown-status.txt")
	if err != nil {
		t.Fatal(err)
	}

	events := ParseJsonsToGoTestEvents(file)

	assert.NotEmpty(t, events)

	events = TrimGoTestEvents(events)

	assert.NotEmpty(t, events)

	containers := ExtractContainers(events)

	results := ExtractResults(events, containers)

	assert.NotEmpty(t, results)

	outputFolder := "allure-results"

	CreateOutputFolder(outputFolder)

	PrintContainers(outputFolder, containers)

	PrintResults(outputFolder, results)

	assert.NotEmpty(t, results)
}
