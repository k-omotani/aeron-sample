#!/bin/bash
# Start Aeron Media Driver
# Requires Java and aeron-all.jar

set -e

# macOS uses /tmp instead of /dev/shm
AERON_DIR=${AERON_DIR:-"/tmp/aeron-default"}
AERON_JAR=${AERON_JAR:-"$HOME/.m2/repository/io/aeron/aeron-all/1.44.1/aeron-all-1.44.1.jar"}

# Check if aeron jar exists
if [ ! -f "$AERON_JAR" ]; then
    echo "Aeron JAR not found at: $AERON_JAR"
    echo "Please download aeron-all.jar or set AERON_JAR environment variable"
    echo ""
    echo "You can download it with Maven:"
    echo "  mvn dependency:get -Dartifact=io.aeron:aeron-all:1.44.1"
    echo ""
    echo "Or download directly:"
    echo "  curl -L -o aeron-all.jar https://repo1.maven.org/maven2/io/aeron/aeron-all/1.44.1/aeron-all-1.44.1.jar"
    exit 1
fi

# Clean previous driver files
echo "Cleaning previous Aeron directory: $AERON_DIR"
rm -rf "$AERON_DIR"

echo "Starting Aeron Media Driver..."
echo "  AERON_DIR: $AERON_DIR"
echo "  AERON_JAR: $AERON_JAR"
echo ""
echo "Press Ctrl+C to stop"
echo ""

export AERON_DIR="$AERON_DIR"

java \
    --add-opens java.base/sun.nio.ch=ALL-UNNAMED \
    --add-opens java.base/java.nio=ALL-UNNAMED \
    --add-opens java.base/java.lang=ALL-UNNAMED \
    --add-opens java.base/jdk.internal.misc=ALL-UNNAMED \
    -cp "$AERON_JAR" \
    io.aeron.driver.MediaDriver
