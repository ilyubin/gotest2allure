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
		bResult, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("error marshal result: %v\n", result)
			continue
		}
		bResult2, err := prettyPrint(bResult)
		if err != nil {
			fmt.Printf("error prettify result: %v\n", bResult)
			continue
		}
		file := path.Join(outputFlag, fmt.Sprintf("%s-result.json", result.UUID))
		err = ioutil.WriteFile(file, bResult2, 0644)
		if err != nil {
			fmt.Printf("error write result: %v to file: %s\n", bResult2, file)
		}
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
