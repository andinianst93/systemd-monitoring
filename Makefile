.PHONY: run build clean

run:
	go run main.go

build:
	mkdir -p bin
	go build -o bin/monitor .

clean:
	rm -rf bin
