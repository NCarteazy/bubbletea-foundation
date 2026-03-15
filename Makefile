.PHONY: test test-race vet build example clean

test:
	go test ./... -v

test-race:
	go test -race ./...

vet:
	go vet ./...

build:
	go build ./...

example:
	cd example && go run .

clean:
	rm -f example/example
