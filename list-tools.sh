#!/bin/bash

echo "AI CLI Manager - Available Tools"
echo "================================="
echo ""
echo "The following 17 AI tools are built into the compiled binary:"
echo ""

cat << 'EOF'
  1. Claude Code           - Anthropic's Claude AI coding assistant
  2. Gemini CLI            - Google's Gemini AI CLI tool
  3. OpenAI Codex          - OpenAI Codex code generation and completion
  4. Qwen CLI              - Alibaba's Qwen AI assistant CLI
  5. GitHub Copilot CLI    - GitHub Copilot command-line interface
  6. Qodo                  - AI test generation and code quality
  7. Ollama                - Run large language models locally
  8. LM Studio CLI         - Local LLM management
  9. Sourcegraph Cody      - Sourcegraph's AI coding assistant
 10. Amazon Q              - Amazon's AI developer assistant
 11. Tabnine CLI           - AI code completion
 12. Pieces CLI            - AI-powered code snippet manager
 13. Mentat                - AI coding assistant with context awareness
 14. GPT Engineer          - AI engineer that builds entire codebases
 15. Auto-GPT              - Autonomous GPT-4 agent
 16. Open Interpreter      - Natural language interface for computers
 17. Sweep AI              - AI-powered code reviewer
EOF

echo ""
echo "Tools with MCP server integration (for Claude Desktop):"
echo "  • Claude Code (filesystem, github)"
echo "  • Ollama (ollama server)"
echo "  • Open Interpreter (code-execution)"

echo ""
echo "To run the application:"
echo "  make run"
echo "  # or"
echo "  ./build/ai-cli-manager"
echo ""
echo "After first run, configuration will be saved to:"
echo "  ~/.ai-cli-manager/tools.json"
echo ""
echo "The tools are embedded in the binary and accessible without external JSON files."