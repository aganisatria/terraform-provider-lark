default: fmt lint install generate

build:
	go build -v ./...

install:
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	go fmt ./...

test:
	go test -gcflags="all=-N -l" -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -gcflags="all=-N -l" -v -cover -timeout 120m ./...

testall:
	make test
	make testacc

.PHONY: fmt lint test testacc build install generate testall
