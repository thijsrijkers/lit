#!/bin/bash

APP_NAME="testapp"                    # Your binary name
SOURCE_FILE="./app/testApp.go"      # Path to the file you just made
OUT_DIR="./base/bin"                 # Output directory
GOOS_TARGET="linux"
GOARCH_TARGET="amd64"

echo "🔨 Building $APP_NAME from $SOURCE_FILE for $GOOS_TARGET/$GOARCH_TARGET..."
mkdir -p "$OUT_DIR"

env GOOS=$GOOS_TARGET GOARCH=$GOARCH_TARGET CGO_ENABLED=0 go build -o "$OUT_DIR/$APP_NAME" "$SOURCE_FILE"

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Built $APP_NAME and copied to $OUT_DIR"
echo "📦 Place this in your lit.yml as:"
echo "    image: \"$APP_NAME\""

