all: build

build:
	go build -o ./bin/app cmd/main.go

build-migrate:
	go build -o ./bin/migrate internal/db/migrate.go

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin

rebuild: clean build

todo:
	echo "---" && grep -rn * -e "TODO:" | grep -v "Makefile:"
