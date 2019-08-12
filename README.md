# gotest2allure
Adapter for `go test` to `allure`


## Install

```bash
go get github.com/ilyubin/gotest2allure/cmd/gotest2allure

```

## Usage

```bash
go test -json > json-report.txt
./gotest2allure -f json-report.txt -o allure-results 
allure serve allure-results
```


```bash
go test -tags e2e -json ./e2e/... > json-report.txt
```

## Inspired by

- https://github.com/tebeka/go2xunit
- https://github.com/jstemmer/go-junit-report
- https://github.com/GabbyyLS/allure-go-common
