package main

import (
	"io"
	"os"
)

func main() {
	piper, pipew := io.Pipe()
	defer piper.Close()

	go func() {
		defer pipew.Close()
		_, _ = io.Copy(pipew, os.Stdin)
	}()

	events := trimGoTestEvents(parseJsonsToGoTestEvents(piper))
	containers := extractContainers(events)
	results := extractResults(events, containers)

	createFolderForAllureResults()
	printContainers(containers)
	printResults(results)
}
