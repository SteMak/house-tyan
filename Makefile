BUILD_DIR=$(shell pwd)/bin
VANILLA_DIR=$(shell pwd)/cli/bot

clean:
	rm -rf $(BUILD_DIR)/*

build: clean
	mkdir -p $(BUILD_DIR)
	cp -R $(VANILLA_DIR)/config $(BUILD_DIR)
	cp -R $(VANILLA_DIR)/assets $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -v -o $(BUILD_DIR)/bot $(VANILLA_DIR)

mod: 
	rm -rf go.mod
	rm -rf go.mod
	go mod init 
	go mod tidy
