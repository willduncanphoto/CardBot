.PHONY: build test clean qa-050

build:
	go build -ldflags="-s -w" -o cardbot .

test:
	go test ./... -count=1 -race

qa-050:
	./scripts/qa_050_smoke.sh

clean:
	rm -f cardbot coverage.out
