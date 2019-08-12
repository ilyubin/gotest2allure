# gotest2allure
Adapter for `go test` to `allure`


## Install

```bash
go get github.com/ilyubin/gotest2allure/cmd/gotest2allure

```

## Usage

Run your tests with flag `-json` and save results to the file `json-report.txt`:

```bash
go test -json > json-report.txt
```

or for e2e tests:

```bash
go test -tags e2e -json ./e2e/... > json-report.txt
```

Run `gotest2allure` from `bin` folder:

```bash
$GOPATH/bin/gotest2allure -f json-report.txt -o allure-results 
```

Generate report with `allure`:

```bash
allure serve allure-results
```

## Inspired by

- https://github.com/tebeka/go2xunit
- https://github.com/jstemmer/go-junit-report
- https://github.com/GabbyyLS/allure-go-common
