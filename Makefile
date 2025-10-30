test:
	go test ./...

tidy:
	go mod tidy


# Variables
BIN:=retroforge
PKG:=./cmd/retroforge
SCALE?=3
CART?=examples/moon-lander.rf

.PHONY: debug release run pack pack-hello clean help

help:
	@echo "make debug        # build $(BIN) with debug info"
	@echo "make release      # build $(BIN) with -s -w"
	@echo "make pack         # pack a cart directory: DIR=<dir> (default examples/moon-lander)"
	@echo "make run          # run a cart: CART=$(CART) (uses -window -scale $(SCALE))"
	@echo "make bundle       # build self-contained binary from CART=<file.rf> OUT=<name>"
	@echo "make wasm         # build WebAssembly binary to webapp/public/engine"
	@echo "make test         # run unit tests"
	@echo "make tidy         # go mod tidy"
	@echo "make clean        # remove binary"

debug:
	go build -o $(BIN) $(PKG)

release:
	go build -ldflags "-s -w" -o $(BIN) $(PKG)

pack:
	@[ -n "$(DIR)" ] || DIR=examples/moon-lander; \
	./$(BIN) -pack $$DIR

run: debug
	./$(BIN) -cart $(CART) -window -scale $(SCALE)

pack-hello: debug
	./$(BIN) -pack examples/helloworld

pack-moon: debug
	./$(BIN) -pack examples/moon-lander

# Build a self-contained executable with embedded cart
bundle: debug
	@[ -n "$(CART)" ] || (echo "CART=<file.rf> required" && false)
	mkdir -p cmd/cartbundle
	cp $(CART) cmd/cartbundle/cart.rf
	go build -o cart-$(shell basename $(CART) .rf) ./cmd/cartbundle
	rm -f cmd/cartbundle/cart.rf

clean:
	rm -f $(BIN)

# Build WASM binary and place alongside wasm_exec.js for the webapp
wasm:
	mkdir -p ../retroforge-webapp/public/engine
	@if [ -f "$$($(shell go env GOROOT))/misc/wasm/wasm_exec.js" ]; then \
	  cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ../retroforge-webapp/public/engine/ ; \
	elif [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
	  cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ../retroforge-webapp/public/engine/ ; \
	else \
	  echo "wasm_exec.js not found; please locate it in your Go installation"; exit 1; \
	fi
	GOOS=js GOARCH=wasm go build -o ../retroforge-webapp/public/engine/retroforge.wasm ./cmd/wasm

.PHONY: web-carts
web-carts:
	mkdir -p ../retroforge-webapp/public/carts
	cp examples/helloworld.rf ../retroforge-webapp/public/carts/helloworld.rf
	cp examples/moon-lander.rf ../retroforge-webapp/public/carts/moon-lander.rf


