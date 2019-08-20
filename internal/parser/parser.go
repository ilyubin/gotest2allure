package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ilyubin/gotest2allure/internal/prefix"
)

const (
	actionRun           = "run"
	actionOutput        = "output"
	resultStatusPassed  = "passed"
	resultStatusFailed  = "failed"
	resultStatusSkipped = "skipped"
)

//ParseJsonsToGoTestEvents ...
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

//TrimGoTestEvents ...
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

//ExtractResults ...
func ExtractResults(events []*GoTestEvent) map[string]*AllureResult {
	results := make(map[string]*AllureResult)
	var isErrorEventContext bool
	var isPanicContext bool
	var isRequestContext bool
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

		if event.Action == actionRun {
			_uuid := getUUID()
			splits := strings.Split(event.Test, "/")

			_time := event.Time.UnixNano() / int64(time.Millisecond)
			result := &AllureResult{
				UUID:      _uuid,
				Name:      event.Test,
				FullName:  event.Test,
				Start:     _time,
				Stop:      _time,
				HistoryID: getUUID(),
				Labels: append(
					resolveSuiteLabels(splits),
					resolvePackageLabel(event),
				),
			}
			results[event.Test] = result
		}

		if event.Action == "pass" {
			result, _ := results[event.Test]
			result.Status = resultStatusPassed
			result.Stop = result.Start + int64(event.Elapsed*1000)
			isPanicContext = false
			isErrorEventContext = false
			isRequestContext = false
			continue
		}
		if event.Action == "fail" {
			result, _ := results[event.Test]
			result.Status = resultStatusFailed
			result.Stop = result.Start + int64(event.Elapsed*1000)
			isPanicContext = false
			isErrorEventContext = false
			isRequestContext = false
			continue
		}
		if event.Action == "skip" {
			result, _ := results[event.Test]
			result.Status = resultStatusSkipped
			result.Stop = result.Start + int64(event.Elapsed*1000)
			isPanicContext = false
			isErrorEventContext = false
			isRequestContext = false

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

			// Handle panic in test
			if strings.HasPrefix(event.Output, "SIGQUIT:") || strings.HasPrefix(event.Output, "panic: runtime error:") && !isPanicContext {
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

			// Handle prefix
			if strings.HasPrefix(output, prefix.Feature) {
				result.Labels = append(result.Labels, Label{
					Name:  "feature",
					Value: strings.Replace(output, prefix.Feature, "", 1),
				})
				continue
			}
			if strings.HasPrefix(output, prefix.Story) {
				result.Labels = append(result.Labels, Label{
					Name:  "story",
					Value: strings.Replace(output, prefix.Story, "", 1),
				})
				continue
			}
			if strings.HasPrefix(output, prefix.Story) {
				result.Description = strings.Replace(output, prefix.Description, "", 1)
				continue
			}

			// Handle error
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

			// Handle curl and grpc_cli
			if strings.HasPrefix(strings.ToLower(output), "response") {
				isRequestContext = false
				result.Steps = append(result.Steps, Step{
					Name:   output,
					Status: "passed",
				})
				continue
			}
			if strings.HasPrefix(strings.ToLower(output), "roundtrip failed") {
				isRequestContext = false
				result.Steps = append(result.Steps, Step{
					Name:   output,
					Status: "failed",
				})
				continue
			}
			if isRequestContext {
				result.Steps[len(result.Steps)-1].Name += strings.TrimSpace(output)
				continue
			}
			if strings.HasPrefix(output, "curl") || strings.HasPrefix(output, "grpc_cli") {
				isRequestContext = true
			}

			// Mark step failed if error
			stepStatus := "passed"
			if strings.Contains(strings.ToLower(output), "error") {
				stepStatus = "failed"
			}

			result.Steps = append(result.Steps, Step{
				Name:   strings.TrimSpace(output),
				Status: stepStatus,
			})
		}
	}
	return results
}

func getUUID() uuid.UUID {
	u, _ := uuid.NewV4()
	return u
}

func resolveSuiteLabels(splits []string) []Label {
	if len(splits) >= 1 {
		return []Label{
			{
				Name:  "parentSuite",
				Value: splits[0],
			},
		}
	}
	//if len(splits) >= 2 {
	//	return []Label{
	//		{
	//			Name:  "parentSuite",
	//			Value: splits[0],
	//		},
	//		{
	//			Name:  "suite",
	//			Value: splits[1],
	//		},
	//	}
	//}
	//if len(splits) > 3 {
	//	return []Label{
	//		{
	//			Name:  "parentSuite",
	//			Value: strings.Join(splits[:len(splits)-3], "/"),
	//		},
	//		{
	//			Name:  "suite",
	//			Value: splits[len(splits)-3],
	//		},
	//		{
	//			Name:  "subSuite",
	//			Value: splits[len(splits)-2],
	//		},
	//	}
	//}

	return []Label{}
}

func resolvePackageLabel(event *GoTestEvent) Label {
	return Label{
		Name:  "package",
		Value: strings.Replace(event.Package, "/", ".", -1),
	}
}
