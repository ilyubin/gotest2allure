package allure

import (
	"testing"

	"github.com/ilyubin/gotest2allure/internal/prefix"
)

// Feature will add BDD-feature in Behaviors
func Feature(t *testing.T, feature string) {
	t.Logf("%s%s", prefix.Feature, feature)
}

// Story will add BDD-story in Behaviors
func Story(t *testing.T, story string) {
	t.Logf("%s%s", prefix.Story, story)
}

// Description will add detailed description for each test
func Description(t *testing.T, description string) {
	t.Logf("%s%s", prefix.Description, description)
}

// Issue will provide a link with a small bug icon
func Issue(t *testing.T, issue string) {
	t.Logf("%s%s", prefix.Issue, issue)
}
