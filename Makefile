.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -rf ${build_dir}

.PHONY: build
build: clean deps vet fmt

.PHONY: update-deps
update-deps:
	go get buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go@latest
	go get github.com/plantoncloud/pulumi-module-golang-commons@latest
