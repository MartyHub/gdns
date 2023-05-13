default: all

all: tidy generate lint test build

build:
	go build -ldflags="-X 'github.com/MartyHub/gdns/version.Version=development'" -race

clean:
	rm -fr .coverage.out gdns

generate:
	go generate ./...

lint:
	go vet ./...
	golangci-lint run

test:
	gotest -coverprofile .coverage.out -race -timeout 10s

tidy:
	go mod tidy
