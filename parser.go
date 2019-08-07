package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

func parseJsonsToGoTestEvents(file io.Reader) []*GoTestEvent {
	reader := bufio.NewReader(file)
	list := make([]*GoTestEvent, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		event := &GoTestEvent{}
		err = json.Unmarshal(line, event)
		if err != nil {
			panic(err)
		}
		list = append(list, event)
	}
	return list
}

func trimGoTestEvents(events []*GoTestEvent) []*GoTestEvent {
	list := make([]*GoTestEvent, 0)
	for _, event := range events {
		if strings.HasSuffix(event.Output, "\n") {
			event.Output = event.Output[:len(event.Output)-1]
		}
		if strings.HasPrefix(event.Output, "    \t\t") {
			list[len(list)-1].Output = list[len(list)-1].Output + event.Output[6:]
			continue
		}
		for i := 0; i < 10; i++ {
			if strings.HasPrefix(event.Output, strings.Repeat(" ", 4*i)+"\t") {
				event.Output = event.Output[(4*i + 1):]
			}
		}
		list = append(list, event)
	}
	return list
}

func extractContainers(events []*GoTestEvent) []*AllureContainer {
	containers := make([]*AllureContainer, 0)
	for _, t2 := range events {
		if t2.Action == "run" && !strings.ContainsAny(t2.Test, "/") {
			container := &AllureContainer{
				UUID:  sUuid(),
				name:  t2.Test,
				//Start: t2.Time.UnixNano() / int64(time.Millisecond),
			}
			containers = append(containers, container)
		}
	}
	return containers
}

func extractResults(events []*GoTestEvent, containers []*AllureContainer) []*AllureResult {
	results := make([]*AllureResult, 0)
	for _, t2 := range events {
		splits := strings.Split(t2.Test, "/")
		//if len(splits) == 1 {
		//	continue
		//}
		if t2.Action == "run" {
		//if t2.Action == "run" && len(splits) == 2 {
			_uuid := sUuid()

			for _, container := range containers {
				if container.name == splits[0] {
					container.Children = append(container.Children, _uuid)
				}
			}

			result := &AllureResult{
				UUID:      _uuid,
				Name:      t2.Test,
				FullName:  t2.Test,
				Start:     t2.Time.UnixNano() / int64(time.Millisecond),
				HistoryID: sUuid(),
				Labels: getLabels(splits),
			}
			results = append(results, result)
		}
	}
	var isErrorEventContext bool
	for _, t3 := range events {
		//splits := strings.Split(t3.Test, "/")
		//if t3.Action == "output" && len(splits) == 2 && !strings.HasPrefix(t3.Output, "===") {
		if t3.Action == "output" && !strings.HasPrefix(t3.Output, "===") {
			for _, result := range results {
				if result.Name == t3.Test {
					if strings.Contains(t3.Output, "--- PASS:") {
						result.Status = "passed"
						continue
					}
					if strings.Contains(t3.Output, "--- FAIL:") {
						result.Status = "failed"
						continue
					}
					if strings.Contains(t3.Output, "--- SKIP:") {
						result.Status = "skipped"
						continue
					}

					reg := regexp.MustCompile(`.+\.go:\d+:\s(.*)`)
					output := reg.ReplaceAllString(t3.Output, "${1}")

					if output == "" {
						continue
					}
					if strings.HasPrefix(output, "Error Trace:") {
						result.StatusDetails.Trace += "\n" + output
						isErrorEventContext = true
						continue
					}
					if strings.HasPrefix(output, "Error:") {
						result.StatusDetails.Message += "\n" + output
						result.StatusDetails.Trace += "\n" + output
						continue
					}
					if isErrorEventContext && !strings.HasPrefix(output, "Test:") {
						result.StatusDetails.Trace += "\n" + output
						continue
					}
					if strings.HasPrefix(output, "Test:") {
						result.StatusDetails.Trace += "\n" + output
						isErrorEventContext = false
						continue
					}

					step := Step{
						Name:   output,
						Status: "passed",
					}

					result.Steps = append(result.Steps, step)
				}
			}
		}
	}
	return results
}

func getLabels(splits []string) []Label {
	if len(splits) == 0 || len(splits) == 1 {
		return []Label{}
	}
	if len(splits) == 2 {
		return []Label{
			{
				Name:  "parentSuite",
				Value: splits[0],
			},
		}
	}
	if len(splits) == 3 {
		return []Label{
			{
				Name:  "parentSuite",
				Value: splits[0],
			},
			{
				Name:  "suite",
				Value: splits[1],
			},
		}
	}
	return []Label{
		{
			Name:  "parentSuite",
			Value: strings.Join(splits[:len(splits)-3], "/"),
		},
		{
			Name:  "suite",
			Value: splits[len(splits)-3],
		},
		{
			Name:  "subSuite",
			Value: splits[len(splits)-2],
		},
	}
}

func sUuid() string {
	uuid4, _ := uuid.NewV4()
	return fmt.Sprintf("%s", uuid4)
}

func createFolderForAllureResults() {
	_ = os.RemoveAll("allure-results")
	_ = os.MkdirAll("allure-results", os.ModePerm)
}

func printResults(results []*AllureResult) {
	for _, result := range results {
		bResult, _ := json.Marshal(result)
		bResult2, _ := prettyPrint(bResult)
		_ = ioutil.WriteFile(fmt.Sprintf("allure-results/%s-result.json", result.UUID), bResult2, 0644)
	}
}

func printContainers(containers []*AllureContainer) {
	for _, container := range containers {
		bContainer, _ := json.Marshal(container)
		bContainer2, _ := prettyPrint(bContainer)
		_ = ioutil.WriteFile(fmt.Sprintf("allure-results/%s-container.json", container.UUID), bContainer2, 0644)
	}
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}
