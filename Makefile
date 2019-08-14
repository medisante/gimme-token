.PHONY: all lint test clean

all:
	go build -o bin/gimme-token .

lint:
	golangci-lint run

test:
	go test -race -cover ./...

clean:
	rm -rf bin
