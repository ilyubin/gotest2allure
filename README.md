# gotest2allure
Adapter for `go test` to `allure`

## Usage

```bash
go test -json > tests.txt
cat tests.txt | ./gotest2allure
allure serve allure-results
```

## Inspired by

- https://github.com/tebeka/go2xunit
- https://github.com/jstemmer/go-junit-report
- https://github.com/GabbyyLS/allure-go-common


## Versions

1.0.0 first stable alpha version
1.0.1 parse result duration