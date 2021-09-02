# Makefile
# @author danpoland

.PHONY: fmt
fmt:
	gofmt -l -s -w . && \
	goimports -l -w . &&  \
	go mod tidy

.PHONY: test
test:
	go test  ./...
