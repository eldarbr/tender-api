all: build

build:
	go build -o ./build/app cmd/main.go

build-migrate:
	go build -o ./build/migrate internal/db/migrate.go

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf build

rebuild: clean build

todo:
	echo "---" && grep -rn * -e "TODO:" | grep -v "Makefile:"
