# Build Release Binary

Build a release-ready Stash binary for the current platform.

## When to use
- Building for production deployment
- Creating distributable binaries
- Testing release builds locally

## Steps

### Local release build

```bash
make build-release
```

This builds dynamically-linked PIE binaries with debug info stripped.

### Custom build flags

```bash
# Release build (stripped + PIE)
make flags-release stash

# Static PIE build (fully self-contained)
make flags-static-pie stash

# Static Windows build
make flags-static-windows stash
```

### macOS .app bundle

```bash
make stash-macapp
```

Creates `Stash.app` in the project root.

### Cross-compilation (requires Docker)

```bash
# Start the compiler container
make start-compiler-container

# Build for specific platforms
docker exec -t build /bin/bash -c "make build-cc-windows"
docker exec -t build /bin/bash -c "make build-cc-linux"
docker exec -t build /bin/bash -c "make build-cc-macos"

# Build all platforms
docker exec -t build /bin/bash -c "make build-cc-all"

# Clean up
make remove-compiler-container
```

Binaries are output to `dist/`.

### Full release pipeline

```bash
make release
```

This runs: `pre-ui -> generate -> ui -> build-release`

## Build metadata

Build info is embedded via ldflags:
- `buildstamp` — build date
- `githash` — short git hash
- `version` — git tag or describe
- `officialBuild` — official build flag

Override with environment variables:
```bash
BUILD_DATE=2024-01-01 GITHASH=abc1234 STASH_VERSION=v0.25.0 OFFICIAL_BUILD=true make stash
```
