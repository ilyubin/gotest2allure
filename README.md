# gotest2allure
Adapter for `go test` to `allure`

## Usage

```bash
go test -json > tests.txt
cat tests.txt | ./gotest2allure
allure serve allure-results
```