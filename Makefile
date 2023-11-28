build:
	go build -o ./bin/mebrox .

run: build
	./bin/mebrox

test:
	go test ./... -v --race