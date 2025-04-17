#!/bin/bash

#
# Make it executable: chmod +x lit-build.sh
#
# ----------------------
# Configurable variables
# ----------------------

APP_NAME="your_app"            # The name of your binary inside the container
SOURCE_DIR="./cmd/app"         # Path to your Go main package
OUT_DIR="./base/bin"           # Where to copy the built binary
GOOS_TARGET="linux"            # Target OS for build
GOARCH_TARGET="amd64"          # Target architecture

# ----------------------
# Start build
# ----------------------

echo "🔨 Building $APP_NAME for $GOOS_TARGET/$GOARCH_TARGET..."
mkdir -p "$OUT_DIR"

# Static build for Linux container
env GOOS=$GOOS_TARGET GOARCH=$GOARCH_TARGET CGO_ENABLED=0 go build -o "$OUT_DIR/$APP_NAME" "$SOURCE_DIR"

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

# ----------------------
# Done!
# ----------------------

echo "✅ Built $APP_NAME and copied to $OUT_DIR"
echo "📦 Place this in your lit.yml as:"
echo "    image: \"$APP_NAME\""

