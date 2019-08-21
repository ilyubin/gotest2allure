package allure

import (
	"testing"

	"github.com/ilyubin/gotest2allure/internal/prefix"
)

// Feature ...
func Feature(t *testing.T, feature string) {
	t.Logf("%s%s", prefix.Feature, feature)
}

// Story ...
func Story(t *testing.T, story string) {
	t.Logf("%s%s", prefix.Story, story)
}

// Description - add detailed description for each test to provide as much context to the report reader as you want
func Description(t *testing.T, description string) {
	t.Logf("%s%s", prefix.Description, description)
}

// Issue ...
func Issue(t *testing.T, issue string) {
	t.Logf("%s%s", prefix.Issue, issue)
}
