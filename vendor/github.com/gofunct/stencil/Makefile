.PHONY: build

build:
	@go get github.com/gofunct/gofs
	@go mod vendor
	@cd ./stencil; go install
	@cd ./stencilbin; go install

