test:
	@# Create dummy cart.rf for cartbundle package (needed for go test)
	@touch cmd/cartbundle/cart.rf && echo "# Dummy cart for testing" > cmd/cartbundle/cart.rf || true
	go test ./...

coverage:
	@# Create dummy cart.rf for cartbundle package (needed for go test)
	@touch cmd/cartbundle/cart.rf && echo "# Dummy cart for testing" > cmd/cartbundle/cart.rf || true
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
	@echo "Coverage report written to coverage.out"
	@echo "Run 'make coverage-html' to view HTML report"

coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage HTML report written to coverage.html"
	@echo "Open coverage.html in your browser to view"

tidy:
	go mod tidy


# Variables
BIN:=retroforge
IMGTOOL_BIN:=imgtool
PKG:=./cmd/retroforge
IMGTOOL_PKG:=./cmd/imgtool
SCALE?=3
CART?=examples/moon-lander.rf
FOLDER?=examples/moon-lander

.PHONY: debug release run run-dev pack pack-hello pack-moon clean help wasm web-carts test coverage coverage-html tidy bundle imgtool imgtool-debug imgtool-release imgtool-clean imgtool-test

help:
	@echo "=== RetroForge Engine ==="
	@echo "make debug        # build $(BIN) with debug info"
	@echo "make release      # build $(BIN) with -s -w"
	@echo "make pack         # pack a cart directory: DIR=<dir> (default examples/moon-lander)"
	@echo "make run          # run a cart: CART=$(CART) (uses -window -scale $(SCALE))"
	@echo "make run-dev      # run cart from folder with hot reload: FOLDER=<dir> (default examples/moon-lander)"
	@echo "make bundle       # build self-contained binary from CART=<file.rf> OUT=<name>"
	@echo "make wasm         # build WebAssembly binary to webapp/public/engine"
	@echo ""
	@echo "=== Image Tool ==="
	@echo "make imgtool         # build $(IMGTOOL_BIN) with debug info"
	@echo "make imgtool-release # build $(IMGTOOL_BIN) with -s -w"
	@echo "make imgtool-test    # run imgtool unit tests"
	@echo "make imgtool-clean   # remove $(IMGTOOL_BIN) binary"
	@echo ""
	@echo "=== Testing & Utilities ==="
	@echo "make test         # run unit tests"
	@echo "make coverage     # run tests with coverage report"
	@echo "make coverage-html # generate HTML coverage report"
	@echo "make tidy         # go mod tidy"
	@echo "make clean        # remove binaries"

debug:
	go build -o $(BIN) $(PKG)

release:
	go build -ldflags "-s -w" -o $(BIN) $(PKG)

pack:
	@[ -n "$(DIR)" ] || DIR=examples/moon-lander; \
	./$(BIN) -pack $$DIR

run: debug
	./$(BIN) -cart $(CART) -window -scale $(SCALE)

run-dev: debug
	@[ -n "$(FOLDER)" ] || FOLDER=examples/moon-lander; \
	./$(BIN) -folder $$FOLDER -window -scale $(SCALE)

# Run a specific example by name (e.g., make example moon-lander)
example:
	@[ -n "$(EXAMPLE)" ] || (echo "Usage: make example EXAMPLE=<name>" && echo "Examples: moon-lander, kitchen-sink, galaxy, helloworld" && false)
	@FOLDER=examples/$(EXAMPLE); \
	if [ ! -d $$FOLDER ]; then \
		echo "Example '$(EXAMPLE)' not found in examples/ directory"; \
		echo "Available examples:"; \
		ls -1 examples/ | grep -v "\.rf$$" | sed 's/^/  - /'; \
		false; \
	fi
	$(MAKE) run-dev FOLDER=$$FOLDER

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
	rm -f $(BIN) $(IMGTOOL_BIN)

# Image Tool Build Targets
# Note: Requires github.com/spf13/cobra (run: go get github.com/spf13/cobra)
imgtool:
	go build -o $(IMGTOOL_BIN) $(IMGTOOL_PKG)

imgtool-debug: imgtool

imgtool-release:
	@if ! go list -m github.com/spf13/cobra >/dev/null 2>&1; then \
		echo "Installing cobra dependency..."; \
		go get github.com/spf13/cobra; \
	fi
	go build -ldflags "-s -w" -o $(IMGTOOL_BIN) $(IMGTOOL_PKG)

imgtool-test:
	go test ./internal/imgtool/... -v

imgtool-clean:
	rm -f $(IMGTOOL_BIN)

# Build WASM binary and place alongside wasm_exec.js for the webapp
wasm:
	mkdir -p ../retroforge-webapp/public/engine
	@GOROOT=$$(go env GOROOT); \
	if [ -f "$$GOROOT/misc/wasm/wasm_exec.js" ]; then \
	  cp "$$GOROOT/misc/wasm/wasm_exec.js" ../retroforge-webapp/public/engine/ ; \
	elif [ -f "$$GOROOT/lib/wasm/wasm_exec.js" ]; then \
	  cp "$$GOROOT/lib/wasm/wasm_exec.js" ../retroforge-webapp/public/engine/ ; \
	else \
	  echo "wasm_exec.js not found; please locate it in your Go installation"; exit 1; \
	fi
	GOOS=js GOARCH=wasm go build -o ../retroforge-webapp/public/engine/retroforge.wasm ./cmd/wasm

.PHONY: web-carts
web-carts:
	mkdir -p ../retroforge-webapp/public/carts
	cp examples/helloworld.rf ../retroforge-webapp/public/carts/helloworld.rf
	cp examples/moon-lander.rf ../retroforge-webapp/public/carts/moon-lander.rf


