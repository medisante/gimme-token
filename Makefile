.PHONY: run lint

build:
	go build -o bin/gimme-token .

lint:
	golangci-lint run --enable-all

clean:
	rm -rf bin
