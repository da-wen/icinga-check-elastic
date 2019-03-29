.PHONY: build
build: build-linux build-mac

.PHONY: build-linux
build-linux:
	mkdir -p build/linux
	GOOS=linux go build -ldflags="-s -w" -o build/linux/icinga_check_elastic

.PHONY: build-mac
build-mac:
	mkdir -p build/darwin
	GOOS=darwin go build -ldflags="-s -w" -o build/darwin/icinga_check_elastic