package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func PrintResults(results []*AllureResult) {
	for _, result := range results {
		bResult, _ := json.Marshal(result)
		bResult2, _ := prettyPrint(bResult)
		_ = ioutil.WriteFile(fmt.Sprintf("allure-results/%s-result.json", result.UUID), bResult2, 0644)
	}
}

func PrintContainers(containers []*AllureContainer) {
	for _, container := range containers {
		bContainer, _ := json.Marshal(container)
		bContainer2, _ := prettyPrint(bContainer)
		_ = ioutil.WriteFile(fmt.Sprintf("allure-results/%s-container.json", container.UUID), bContainer2, 0644)
	}
}

func PrintAttachment(attachment Attachment, output string) {
	bOutput := []byte(output)
	_ = ioutil.WriteFile(fmt.Sprintf("allure-results/%s", attachment.Source), bOutput, 0644)
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}
