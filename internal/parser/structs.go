package parser

import (
	"time"

	"github.com/gofrs/uuid"
)

// =====================================================================================================================

// AllureResult ...
type AllureResult struct {
	Name          string       `json:"name"`
	Status        string       `json:"status"`
	StatusDetails StatusDetail `json:"statusDetails"`
	Steps         []Step       `json:"steps"`
	Start         int64        `json:"start"`
	Stop          int64        `json:"stop"`
	UUID          uuid.UUID    `json:"uuid"`
	HistoryID     uuid.UUID    `json:"historyId"`
	FullName      string       `json:"fullName"`
	Labels        []Label      `json:"labels"`
	Description   string       `json:"description"`
	Links         []Link       `json:"links"`
}

// StatusDetail ...
type StatusDetail struct {
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

// Step ...
type Step struct {
	Name        string       `json:"name"`
	Status      string       `json:"status"`
	Attachments []Attachment `json:"attachments"`
	Start       int64        `json:"start"`
	Stop        int64        `json:"stop"`
	Parameters  []Parameter  `json:"parameters"`
}

// Attachment ...
type Attachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Data   string `json:"-"`
}

// Label ...
type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Parameter ...
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Link ...
type Link struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// =====================================================================================================================

//GoTestEvent - line in json-report.txt file
type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"`
	Output  string    `json:"Output"`
}
