#!/bin/bash
set -e

# Aeron directory in shared memory
AERON_DIR=/dev/shm/aeron

# Clean up any existing Aeron directory
rm -rf "$AERON_DIR" 2>/dev/null || true
mkdir -p "$AERON_DIR"

echo "Starting Aeron Media Driver..."
echo "AERON_DIR: $AERON_DIR"

exec java \
    --add-opens java.base/sun.nio.ch=ALL-UNNAMED \
    --add-opens java.base/java.nio=ALL-UNNAMED \
    --add-opens java.base/java.lang=ALL-UNNAMED \
    --add-opens java.base/jdk.internal.misc=ALL-UNNAMED \
    -Daeron.dir="$AERON_DIR" \
    -Daeron.mtu.length=1408 \
    -Daeron.threading.mode=SHARED \
    -cp /opt/aeron-all.jar \
    io.aeron.driver.MediaDriver
