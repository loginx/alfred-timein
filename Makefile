.PHONY: all test build clean alfredworkflow

BIN_DIR := bin

all: build

test:
	go test ./...

build:
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o $(BIN_DIR)/geotz_amd64 ./cmd/geotz
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o $(BIN_DIR)/geotz_arm64 ./cmd/geotz
	lipo -create -output $(BIN_DIR)/geotz $(BIN_DIR)/geotz_amd64 $(BIN_DIR)/geotz_arm64
	rm $(BIN_DIR)/geotz_amd64 $(BIN_DIR)/geotz_arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o $(BIN_DIR)/timein_amd64 ./cmd/timein
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o $(BIN_DIR)/timein_arm64 ./cmd/timein
	lipo -create -output $(BIN_DIR)/timein $(BIN_DIR)/timein_amd64 $(BIN_DIR)/timein_arm64
	rm $(BIN_DIR)/timein_amd64 $(BIN_DIR)/timein_arm64

alfredworkflow: build
	cp $(BIN_DIR)/geotz $(BIN_DIR)/timein workflow/
	cd workflow && zip -r ../TimeIn.alfredworkflow . -x '*.DS_Store'
	rm workflow/geotz workflow/timein

clean:
	rm -rf $(BIN_DIR)/*.alfredworkflow $(BIN_DIR)/geotz $(BIN_DIR)/timein 