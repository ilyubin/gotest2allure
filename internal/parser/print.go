package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func CreateOutputFolder(folder string) {
	_ = os.RemoveAll(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
}

func PrintResults(outputFlag string, results map[string]*AllureResult) {
	for _, result := range results {
		bResult, _ := json.Marshal(result)
		bResult2, _ := prettyPrint(bResult)
		_ = ioutil.WriteFile(path.Join(outputFlag, fmt.Sprintf("%s-result.json", result.UUID)), bResult2, 0644)
	}
}

func PrintContainers(outputFlag string, containers []*AllureContainer) {
	for _, container := range containers {
		bContainer, _ := json.Marshal(container)
		bContainer2, _ := prettyPrint(bContainer)
		_ = ioutil.WriteFile(path.Join(outputFlag, fmt.Sprintf("%s-container.json", container.UUID)), bContainer2, 0644)
	}
}

func PrintAttachment(outputFlag string, attachment Attachment, output string) {
	bOutput := []byte(output)
	_ = ioutil.WriteFile(path.Join(outputFlag, fmt.Sprintf("%s", attachment.Source)), bOutput, 0644)
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}
