BINARY  := paperflow
MODULE  := github.com/fcarvajalbrown/agentic-cs-paper-makers
CMD     := ./cmd/paperflow

.PHONY: build test clean cross

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY) $(BINARY).exe

cross:
	GOOS=linux   GOARCH=amd64  go build -o dist/$(BINARY)-linux-amd64   $(CMD)
	GOOS=linux   GOARCH=arm64  go build -o dist/$(BINARY)-linux-arm64   $(CMD)
	GOOS=darwin  GOARCH=amd64  go build -o dist/$(BINARY)-darwin-amd64  $(CMD)
	GOOS=darwin  GOARCH=arm64  go build -o dist/$(BINARY)-darwin-arm64  $(CMD)
	GOOS=windows GOARCH=amd64  go build -o dist/$(BINARY)-windows-amd64.exe $(CMD)
