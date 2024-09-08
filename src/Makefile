all: build

build:
	go build -o ./build/ cmd/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf build

rebuild: clean build

todo:
	echo "---" && grep -rn * -e "TODO:" | grep -v "Makefile:"
