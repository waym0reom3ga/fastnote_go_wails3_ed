.PHONY: build test clean

BIN := fastnote-wails3

build:
	PATH="/home/waymore/go/bin:$(PATH)" wails3 build

test: build
	./bin/fastnote-wails3 --version

clean:
	rm -f $(BIN)
	rm -rf bin/ build/linux/appimage/
