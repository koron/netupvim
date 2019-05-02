DESCRIBE_TAG = $(shell git describe --tags)
VERSION = $(DESCRIBE_TAG)

default: test

test:
	go test ./...

lint:
	go vet ./...
	@echo ""
	golint ./...

cyclo:
	-gocyclo -top 10 -avg .

report:
	@echo "misspell"
	@find . -name \*.go | xargs misspell
	@echo ""
	-gocyclo -over 14 -avg .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

deps-list:
	@go list -f '{{join .Imports "\n"}}' ./... | sort -u | grep -v `go list`

deps-update:
	@go list -f '{{join .Imports "\n"}}' ./... | sort -u | grep -v `go list` | xargs go get -u -d -v

exe: netupvim.exe

release: clean zip
.PHONY: release

zip: netupvim-$(VERSION).zip

clean:
	go clean
	rm -f netupvim-v*.zip

.PHONY: test lint cyclo report zip

netupvim-$(VERSION).zip: netupvim.exe
	zip -r9 netupvim-$(VERSION).zip netupvim.exe UPDATE.bat RESTORE.bat

netupvim.exe:
	GOOS=windows GOARCH=386 go build -ldflags="-X main.version=$(VERSION)"
