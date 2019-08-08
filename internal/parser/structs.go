package parser

import (
	"time"
)

// =====================================================================================================================

type AllureResult struct {
	Name          string       `json:"name"`
	Status        string       `json:"status"`
	StatusDetails StatusDetail `json:"statusDetails"`
	Steps         []Step       `json:"steps"`
	Start         int64        `json:"start"`
	Stop          int64        `json:"stop"`
	UUID          string       `json:"uuid"`
	HistoryID     string       `json:"historyId"`
	FullName      string       `json:"fullName"`
	Labels        []Label      `json:"labels"`
}

type StatusDetail struct {
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

type Step struct {
	Name        string       `json:"name"`
	Status      string       `json:"status"`
	Attachments []Attachment `json:"attachments"`
	Start       int64        `json:"start"`
	Stop        int64        `json:"stop"`
}

type Attachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Data   string `json:"-"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// =====================================================================================================================

type AllureContainer struct {
	UUID     string   `json:"uuid"`
	Children []string `json:"children"`
	Befores  []Before `json:"befores"`
	Afters   []After  `json:"afters"`
	Start    int64    `json:"start"`
	Stop     int64    `json:"stop"`
	name     string
}

type Before struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Start  int64  `json:"start"`
	Stop   int64  `json:"stop"`
}

type After struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Start  int64  `json:"start"`
	Stop   int64  `json:"stop"`
}

// =====================================================================================================================

type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"`
	Output  string    `json:"Output"`
}
