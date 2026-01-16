#!/bin/bash
# Start Aeron Archive Media Driver
# Requires Java and aeron-all.jar

set -e

# macOS uses /tmp instead of /dev/shm
AERON_DIR=${AERON_DIR:-"/tmp/aeron-default"}
ARCHIVE_DIR=${ARCHIVE_DIR:-"/tmp/aeron-archive"}
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
echo "Cleaning previous Aeron directories..."
rm -rf "$AERON_DIR"
rm -rf "$ARCHIVE_DIR"
mkdir -p "$ARCHIVE_DIR"

echo "Starting Aeron Archive Media Driver..."
echo "  AERON_DIR: $AERON_DIR"
echo "  ARCHIVE_DIR: $ARCHIVE_DIR"
echo "  AERON_JAR: $AERON_JAR"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Export as environment variables (Aeron reads these)
export AERON_DIR="$AERON_DIR"
export AERON_ARCHIVE_DIR="$ARCHIVE_DIR"
export AERON_ARCHIVE_CONTROL_CHANNEL="aeron:udp?endpoint=localhost:8010"
export AERON_ARCHIVE_CONTROL_STREAM_ID=10
export AERON_ARCHIVE_CONTROL_RESPONSE_CHANNEL="aeron:udp?endpoint=localhost:0"
export AERON_ARCHIVE_CONTROL_RESPONSE_STREAM_ID=20
export AERON_ARCHIVE_RECORDING_EVENTS_CHANNEL="aeron:udp?endpoint=localhost:8030"
export AERON_ARCHIVE_RECORDING_EVENTS_ENABLED=true

java \
    --add-opens java.base/sun.nio.ch=ALL-UNNAMED \
    --add-opens java.base/java.nio=ALL-UNNAMED \
    --add-opens java.base/java.lang=ALL-UNNAMED \
    --add-opens java.base/jdk.internal.misc=ALL-UNNAMED \
    -cp "$AERON_JAR" \
    io.aeron.archive.ArchivingMediaDriver
