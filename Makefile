
BIN = 4bod-asm
SRC_DIR = src
BIN_DIR = bin
LIN_DIR = $(BIN_DIR)/lin
WIN_DIR = $(BIN_DIR)/win

build: build-lin build-win

build-lin:
	@echo '--> Builing Linux binary...'
	GOOS=linux go build -o $(LIN_DIR)/$(BIN) $(SRC_DIR)/main.go

build-win:
	@echo '--> Builing Windows binary...'
	GOOS=windows go build -o $(LIN_DIR)/$(BIN) $(SRC_DIR)/main.go
