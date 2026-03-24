#!/bin/bash
# Build script for david CalDAV Server

set -e

VERSION="${VERSION:-1.0.0}"
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

OUTPUT_DIR="${OUTPUT_DIR:-./build}"
BINARY_NAME="david"

echo "=== Building david CalDAV Server ==="
echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"
echo "Commit: $GIT_COMMIT"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build binary
echo "Building binary..."
go build -ldflags "-X main.version=$VERSION -X main.buildTime=$BUILD_TIME -X main.gitCommit=$GIT_COMMIT" \
  -o "$OUTPUT_DIR/$BINARY_NAME" \
  ./cmd/david

# Set executable permission
chmod +x "$OUTPUT_DIR/$BINARY_NAME"

echo "✓ Binary built: $OUTPUT_DIR/$BINARY_NAME"

# Generate checksum
echo ""
echo "Generating checksum..."
sha256sum "$OUTPUT_DIR/$BINARY_NAME" | sed "s|$OUTPUT_DIR/||" > "$OUTPUT_DIR/$BINARY_NAME.sha256"
echo "✓ Checksum: $OUTPUT_DIR/$BINARY_NAME.sha256"

# Create release archive
echo ""
echo "Creating release archive..."
ARCHIVE_NAME="$BINARY_NAME-$VERSION-linux-amd64"
tar -czf "$OUTPUT_DIR/$ARCHIVE_NAME.tar.gz" -C "$OUTPUT_DIR" "$BINARY_NAME" "$BINARY_NAME.sha256"
echo "✓ Archive: $OUTPUT_DIR/$ARCHIVE_NAME.tar.gz"

# Generate release notes
echo ""
echo "Release $VERSION:"
echo "- Full CalDAV support (RFC 4791)"
echo "- REST API with multiple auth methods"
echo "- Admin management features"
echo "- IANA timezone support"
echo "- Calendar sharing and collaboration"
echo "- Rate limiting and audit logging"

echo ""
echo "=== Build Complete ==="
echo "Artifacts in: $OUTPUT_DIR/"
ls -la "$OUTPUT_DIR/"
