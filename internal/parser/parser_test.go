package parser

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullFlow(t *testing.T) {
	t.SkipNow()
	//file, err := os.Open("../../testdata/unknown-status.txt")
	//file, err := os.Open("../../testdata/panic3.txt")
	//file, err := os.Open("../../testdata/panic.txt")

	events := data2TestEvents(t, "few-errors.txt")

	assert.NotEmpty(t, events)

	events = TrimGoTestEvents(events)

	assert.NotEmpty(t, events)

	results := ExtractResults(events, "%s", "%s")

	assert.NotEmpty(t, results)

	outputFolder := "allure-results"

	CreateOutputFolder(outputFolder)

	PrintResults(outputFolder, results)

	assert.NotEmpty(t, results)
}

func TestWithTestData(t *testing.T) {
	t.Parallel()

	t.Run("panic", func(t *testing.T) {
		t.Parallel()
		events := data2TestEvents(t, "panic.txt")
		results := ExtractResults(events, "%s", "%s")
		for k, r := range results {
			if k == "TestDesktopWidgets/widget_Characteristics/top_and_long_chatacteristics" {
				assert.Equal(t, "failed", r.Status)
				assert.Contains(t, r.StatusDetails.Message, "panic: runtime error: index out of range")
			}
		}
	})

	t.Run("few errors", func(t *testing.T) {
		t.Parallel()
		events := data2TestEvents(t, "few-errors.txt")
		results := ExtractResults(events, "%s", "%s")
		for k, r := range results {
			if k == "TestWidgetSKSCart/not_credit_item" {
				assert.Equal(t, "skipped", r.Status)
				assert.Equal(t, "GetNotCreditItem", r.StatusDetails.Message)
			}
			if k == "TestWidgetTotalCredit/item_with_credit" {
				assert.Equal(t, "failed", r.Status)
				assert.Equal(t, "Few errors", r.StatusDetails.Message)
			}
		}
	})
}

func data2TestEvents(t *testing.T, fileName string) []*GoTestEvent {
	file, err := os.Open(fmt.Sprintf("../../testdata/%s", fileName))
	require.NoError(t, err)
	events := ParseJsonsToGoTestEvents(file)
	assert.NotEmpty(t, events)
	events = TrimGoTestEvents(events)
	assert.NotEmpty(t, events)
	return events
}
