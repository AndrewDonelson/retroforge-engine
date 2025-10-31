# Cart Bundle

This directory contains the cartbundle command which creates self-contained executables with embedded cart files.

## Important Note: Linter Error (False Positive)

The `cart.rf` file in this directory is a placeholder required for the `//go:embed` directive. It will be overwritten during the bundle build process (see Makefile `bundle` target).

**If you see a linter error about `cart.rf: no matching files found`**: This is a **false positive** from the Go language server (gopls). The file exists at `cmd/cartbundle/cart.rf` and the code compiles successfully. The error can be safely ignored.

### Why this happens:
- The file has build constraint `//go:build !js && !wasm`
- gopls may analyze with different build tags
- The embed directive is checked statically before build

### Verification:
```bash
# File exists
ls -lh cmd/cartbundle/cart.rf

# Code compiles
go build ./cmd/cartbundle
```

The placeholder file can be any valid `.rf` cart file - it gets replaced during build time.

