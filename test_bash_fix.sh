#!/bin/sh
# Test script to verify bash is available in the container

echo "Testing bash availability in PIXELZX container..."

# Test if bash is available
docker run --rm yuchanshin/pixelzx-evm:latest which bash

if [ $? -eq 0 ]; then
    echo "✅ SUCCESS: Bash is available in the container"
else
    echo "❌ FAILED: Bash is not available in the container"
    exit 1
fi

echo "Testing sh compatibility..."
# Test if our scripts work with sh
docker run --rm -v $(pwd)/scripts:/scripts yuchanshin/pixelzx-evm:latest sh -c "sh /scripts/detect-platform.sh --help"

if [ $? -eq 0 ]; then
    echo "✅ SUCCESS: Scripts are sh-compatible"
else
    echo "❌ FAILED: Scripts are not sh-compatible"
    exit 1
fi

echo "✅ All tests passed!"