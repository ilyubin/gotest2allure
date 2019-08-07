package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	file, err := os.Open("testdata/04-cases.txt")
	if err != nil {
		t.Fatal(err)
	}

	events := parseJsonsToGoTestEvents(file)

	assert.NotEmpty(t, events)

	events = trimGoTestEvents(events)

	assert.NotEmpty(t, events)

	containers := extractContainers(events)

	results := extractResults(events, containers)

	assert.NotEmpty(t, results)

	createFolderForAllureResults()

	printContainers(containers)

	printResults(results)

	assert.NotEmpty(t, results)
}
