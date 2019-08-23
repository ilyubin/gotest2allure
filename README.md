[![Go Report Card](https://goreportcard.com/badge/github.com/ilyubin/gotest2allure)](https://goreportcard.com/report/github.com/ilyubin/gotest2allure)
[![Travis](https://travis-ci.org/ilyubin/gotest2allure.svg?branch=master)](https://travis-ci.org/ilyubin/gotest2allure)

# gotest2allure
Covert `go test` results to `allure`


## Install

```bash
go get github.com/ilyubin/gotest2allure/cmd/gotest2allure

```

## Usage

Run your tests with flag `-json` and save results to the file `json-report.txt`:

```bash
go test -json > json-report.txt
// OR
go test -tags e2e -json ./e2e/... > json-report.txt
```

Run `gotest2allure` from `bin` folder:

```bash
$GOPATH/bin/gotest2allure -f json-report.txt 
```

Generate report with `allure`:

```bash
allure serve allure-results
```

## Parameters

* `-o` specify output directory, default `allure-results`
* `-issuePattern` specify tracker pattern, should have substring `%s`

## Allure features

```go
import "github.com/ilyubin/gotest2allure/pkg/allure"
```

### Issue

```go
allure.Issue(t, "TASK-1")
```

```bash
$GOPATH/bin/gotest2allure -f json-report.txt -issuePattern https://my.jira.com/browse/%s
```

### Description

```go
allure.Description(t, "My detailed description")
```

## Inspired by

- https://github.com/tebeka/go2xunit
- https://github.com/jstemmer/go-junit-report
- https://github.com/GabbyyLS/allure-go-common
