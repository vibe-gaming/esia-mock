LOCAL_BIN:=$(CURDIR)/bin
BUILD_DIR:=$(CURDIR)/cmd

# build binary
build: deps build_binary

build-binary:
	@echo 'build backend binary'
	go build -o $(LOCAL_BIN) $(BUILD_DIR)

deps:
	@echo 'install dependencies'
	go mod tidy -v

# run app
run: deps run-app

run-app:
	@echo 'run backend'
	go run $(BUILD_DIR)/app/main.go
