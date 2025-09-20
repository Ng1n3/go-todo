build_dir = ./bin
binary_name = myapp
home_path = ./cmd/

.DEFAULT_GOAL := build

.PHONY:help fmt vet build confirm clean
## help: show help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## fmt: Run go fmt on all packages
fmt:
	go fmt ./...

## vet: Run go vet on all packages
vet: fmt
	go vet ./...

## build: build the application
build: vet
	@mkdir -p ${build_dir}
	GOARCH=amd64 GOOS=linux go build -o ${build_dir}/${binary_name}-linux ${home_path}

confirm: build
	@echo -n 'Are you sure [y/N]' && read ans && [ "$$ans" = "y" ]

## clean: clean up the binary
clean: confirm
	@echo -n "Cleaning up..."
	@rm -rf ${build_dir}
	@echo -n "Removing storage and json files..."
	@rm -rf storage save_todos.json
