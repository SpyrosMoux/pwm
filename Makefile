.DEFAULT_GOAL := build

BIN := bin
APP := pwm
PKG := .          # or ./cmd/pwm if your main package lives there

.PHONY: bin
bin:
	mkdir -p $(BIN)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build-linux
build-linux: bin
	GOOS=linux GOARCH=amd64 go build -o $(BIN)/$(APP)-linux-amd64 $(PKG)

.PHONY: build-linux-arm
build-linux-arm: bin
	GOOS=linux GOARCH=arm64 go build -o $(BIN)/$(APP)-linux-arm64 $(PKG)

.PHONY: build-mac-intel
build-mac-intel: bin
	GOOS=darwin GOARCH=amd64 go build -o $(BIN)/$(APP)-darwin-amd64 $(PKG)

.PHONY: build-mac-arm
build-mac-arm: bin
	GOOS=darwin GOARCH=arm64 go build -o $(BIN)/$(APP)-darwin-arm64 $(PKG)

.PHONY: build-windows
build-windows: bin
	GOOS=windows GOARCH=amd64 go build -o $(BIN)/$(APP)-windows-amd64.exe $(PKG)

.PHONY: build-windows-arm
build-windows-arm: bin
	GOOS=windows GOARCH=arm64 go build -o $(BIN)/$(APP)-windows-arm64.exe $(PKG)

.PHONY: build
build: build-linux build-linux-arm build-mac-intel build-mac-arm build-windows build-windows-arm
