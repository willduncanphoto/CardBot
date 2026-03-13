.PHONY: build test clean

build:
	go build -ldflags="-s -w" -o cardbot .

test:
	go test ./... -count=1 -race

clean:
	rm -f cardbot coverage.out
