gen_code:
	`go env GOPATH`/bin/goctl api format -dir ginctl_gen
	`go env GOPATH`/bin/ginctl -a ginctl_gen/financial_statement.api -d ./src

build:
	rm -f financial_statement
	go mod tidy
	go build src/financial_statement.go

test: 
	go test -coverpkg ./src/logic/user ./src/logic/user

.PHONY: build gen_code test