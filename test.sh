#!/bin/bash

echo "Testing AI CLI Manager build..."

# Check if build was successful
if [ -f "build/ai-cli-manager" ]; then
    echo "✓ Build successful: build/ai-cli-manager exists"
else
    echo "✗ Build failed: build/ai-cli-manager not found"
    exit 1
fi

# Check file size
SIZE=$(du -h build/ai-cli-manager | cut -f1)
echo "✓ Binary size: $SIZE"

# Check if executable
if [ -x "build/ai-cli-manager" ]; then
    echo "✓ Binary is executable"
else
    echo "✗ Binary is not executable"
    exit 1
fi

# List source files
echo ""
echo "Source files created:"
find src -name "*.go" | while read file; do
    echo "  - $file"
done

echo ""
echo "Configuration will be stored in:"
echo "  - ~/.ai-cli-manager/tools.json (tool definitions)"
echo "  - ~/.ai-cli-manager/config.json (GitHub settings)"
echo "  - ~/Library/Application Support/Claude/claude_desktop_config.json (MCP servers)"

echo ""
echo "To run the application in a terminal:"
echo "  make run"
echo "  # or"
echo "  ./build/ai-cli-manager"

echo ""
echo "✓ All checks passed!"