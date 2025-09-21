
PROJECT_NAME=tried
VERSION=0.1.0

GO_SOURCES=$(shell find . -name '*.go' -not -path "./vendor/*")

$(PROJECT_NAME): $(GO_SOURCES)
	go build -o $@ ./cmd/main.go

clean:
	rm -rf $(PROJECT_NAME)

test:
	go test -coverprofile=coverage.out ./pkg/...
	go tool cover -func=coverage.out

all: $(PROJECT_NAME)
