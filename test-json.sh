#!/bin/bash

echo "Testing JSON-based tool loading..."
echo ""

# Check if JSON file exists and is valid
if [ ! -f "ai_tools.json" ]; then
    echo "❌ ai_tools.json not found"
    exit 1
fi

echo "✅ ai_tools.json found"

# Check JSON validity
if jq empty ai_tools.json; then
    echo "✅ JSON is valid"
else
    echo "❌ JSON is invalid"
    exit 1
fi

# Count tools
TOOL_COUNT=$(jq length ai_tools.json)
echo "✅ Found $TOOL_COUNT tools in JSON"

# Show first 5 tool names
echo ""
echo "First 5 tools:"
jq -r '.[0:5][] | "  • " + .name' ai_tools.json

# Show tools with MCP servers
echo ""
echo "Tools with MCP servers:"
jq -r '.[] | select(.mcp_servers != null) | "  • " + .name + " (" + (.mcp_servers | length | tostring) + " servers)"' ai_tools.json

echo ""
echo "✅ JSON configuration is ready!"
echo ""
echo "The application will load these tools from ai_tools.json when run."