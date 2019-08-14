package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	actionRun    = "run"
	actionOutput = "output"
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

		if t2.Action == actionRun && !strings.ContainsAny(t2.Test, "/") {
			container := &AllureContainer{
				UUID: uuid.NewV4(),
				name: t2.Test,
			}
			containers = append(containers, container)
		}
	}
	return containers
}

func ExtractResults(events []*GoTestEvent, containers []*AllureContainer) map[string]*AllureResult {
	results := make(map[string]*AllureResult)
	for _, event := range events {
		splits := strings.Split(event.Test, "/")
		if event.Action == actionRun {
			_uuid := uuid.NewV4()

			for _, container := range containers {
				if container.name == splits[0] {
					container.Children = append(container.Children, _uuid)
				}
			}

			result := &AllureResult{
				UUID:      _uuid,
				Name:      event.Test,
				FullName:  event.Test,
				Start:     event.Time.UnixNano() / int64(time.Millisecond),
				HistoryID: uuid.NewV4(),
				Labels:    getLabels(splits),
			}
			results[event.Test] = result
		}
	}
	var isErrorEventContext bool
	var isPanicContext bool
	for i, event := range events {
		if event.Test == "" {
			continue
		}
		if strings.HasPrefix(event.Output, "===") {
			continue
		}
		if strings.Contains(event.Output, "--- ") {
			continue
		}

		if event.Action == "pass" {
			result, _ := results[event.Test]
			result.Status = "passed"
			result.Stop = result.Start + int64(event.Elapsed*1000)
			isPanicContext = false
			isErrorEventContext = false
			continue
		}
		if event.Action == "fail" {
			result, _ := results[event.Test]
			result.Status = "failed"
			result.Stop = result.Start + int64(event.Elapsed*1000)
			isPanicContext = false
			isErrorEventContext = false
			continue
		}
		if event.Action == "skip" {
			result, _ := results[event.Test]
			result.Status = "skipped"
			result.Stop = result.Start + int64(event.Elapsed*1000)
			isPanicContext = false
			isErrorEventContext = false

			prev := events[i-1]
			if prev.Action != actionOutput {
				continue
			}
			reg := regexp.MustCompile(`(.+_test\.go:\d+):\s(.*)`)
			trace := reg.ReplaceAllString(prev.Output, "${1}")
			output := reg.ReplaceAllString(prev.Output, "${2}")
			if output == "" {
				continue
			}
			if strings.Contains(output, "--- SKIP:") {
				continue
			}
			result.StatusDetails.Message = output
			result.StatusDetails.Trace = trace
			result.Steps = result.Steps[:len(result.Steps)-1]

			continue
		}

		if event.Action == actionOutput {
			result, ok := results[event.Test]
			if !ok {
				fmt.Printf("unexpected event: %v", event)
				continue
			}

			// Panic in test
			if strings.HasPrefix(event.Output, "SIGQUIT:") {
				result.StatusDetails.Message += "\n" + event.Output
				result.StatusDetails.Trace += "\n" + event.Output
				isPanicContext = true
				continue
			}
			if isPanicContext {
				result.StatusDetails.Trace += "\n" + event.Output
				continue
			}

			reg := regexp.MustCompile(`.+\.go:\d+:\s(.*)`)
			output := reg.ReplaceAllString(event.Output, "${1}")

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
			if isErrorEventContext {
				result.StatusDetails.Trace += "\n" + output
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
