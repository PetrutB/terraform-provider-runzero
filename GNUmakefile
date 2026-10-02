default: build

build:
	go build -o terraform-provider-runzero

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/petrutb/runzero/0.1.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-runzero ~/.terraform.d/plugins/registry.terraform.io/petrutb/runzero/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/

test:
	go test ./... -v

testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

lint:
	golangci-lint run ./...

generate:
	go generate ./...

.PHONY: build install test testacc lint generate
