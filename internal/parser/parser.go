package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	run = "run"
)

func ParseJsonsToGoTestEvents(file io.Reader) []*GoTestEvent {
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

func TrimGoTestEvents(events []*GoTestEvent) []*GoTestEvent {
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

func ExtractContainers(events []*GoTestEvent) []*AllureContainer {
	containers := make([]*AllureContainer, 0)
	for _, t2 := range events {

		if t2.Action == run && !strings.ContainsAny(t2.Test, "/") {
			container := &AllureContainer{
				UUID: sUUID(),
				name: t2.Test,
				//Start: t2.Time.UnixNano() / int64(time.Millisecond),
			}
			containers = append(containers, container)
		}
	}
	return containers
}

func ExtractResults(events []*GoTestEvent, containers []*AllureContainer) []*AllureResult {
	results := make([]*AllureResult, 0)
	for _, t2 := range events {
		splits := strings.Split(t2.Test, "/")
		//if len(splits) == 1 {
		//	continue
		//}
		if t2.Action == run {
			//if t2.Action == "run" && len(splits) == 2 {
			_uuid := sUUID()

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
				HistoryID: sUUID(),
				Labels:    getLabels(splits),
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
						result.Stop = result.Start + elapsedMilliSeconds(t3.Output)
						continue
					}
					if strings.Contains(t3.Output, "--- FAIL:") {
						result.Status = "failed"
						result.Stop = result.Start + elapsedMilliSeconds(t3.Output)
						continue
					}
					if strings.Contains(t3.Output, "--- SKIP:") {
						result.Status = "skipped"
						result.Stop = result.Start + 1
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
					//if strings.HasPrefix(output, "curl") || strings.HasPrefix(output, "grpc_cli") {
					//	attachment := Attachment{
					//		Name:   "curl",
					//		Source: sUUID() + "-attachment.txt",
					//		Type:   "text/plain",
					//	}
					//	//step.Attachments = append(step.Attachments, attachment)
					//	//printAttachment(attachment, output)
					//}

					result.Steps = append(result.Steps, step)
				}
			}
		}
	}
	return results
}

func elapsedMilliSeconds(output string) int64 {
	regexStatus := regexp.MustCompile(`--- (PASS|FAIL|SKIP): (.+) \((\d+\.\d+)(?: seconds|s)\)`)
	matches := regexStatus.FindStringSubmatch(output)
	elapsed := parseSeconds(matches[3])
	return elapsed.Nanoseconds() / 1000000
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

func sUUID() string {
	uuid4 := uuid.NewV4()
	return fmt.Sprintf("%s", uuid4)
}

func CreateOutputFolder(folder string) {
	_ = os.RemoveAll(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
}

func parseSeconds(t string) time.Duration {
	if t == "" {
		return time.Duration(0)
	}
	d, _ := time.ParseDuration(t + "s")
	return d
}
