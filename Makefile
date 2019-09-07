BUILD_BIN_NAME=pm-zero
BUILD_DIR=dist

COMMIT_ID := $(shell git rev-parse HEAD)

build:
	# Install build dependencies
	go get github.com/inconshreveable/mousetrap
	rm -Rf $(BUILD_DIR)
	# Linux x86
	env GOOS=linux 	 GOARCH=386 go build -o $(BUILD_DIR)/linux/$(BUILD_BIN_NAME) app.go
	# MacOS X x86
	env GOOS=darwin  GOARCH=386 go build -o $(BUILD_DIR)/macosx/$(BUILD_BIN_NAME) app.go
	# Windows x86
	env GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/windows/$(BUILD_BIN_NAME).exe app.go
