.PHONY: build test selftest clean

BIN := fastnote-wails3

build:
	PATH="/home/waymore/go/bin:$(PATH)" wails3 build

selftest: build
	./bin/fastnote-wails3 --selftest

test: selftest

clean:
	rm -f $(BIN)
	rm -rf bin/ build/linux/appimage/
