package allure

import (
	"testing"

	"github.com/ilyubin/gotest2allure/internal/prefix"
)

func Feature(t *testing.T, feature string) {
	t.Logf("%s%s", prefix.Feature, feature)
}

func Story(t *testing.T, story string) {
	t.Logf("%s%s", prefix.Story, story)
}

func Description(t *testing.T, description string) {
	t.Logf("%s%s", prefix.Description, description)
}
