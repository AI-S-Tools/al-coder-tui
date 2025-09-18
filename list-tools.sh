#!/bin/bash

echo "AI CLI Manager - Available Tools"
echo "================================="
echo ""
echo "The following 20 AI tools are built into the compiled binary:"
echo ""

cat << 'EOF'
  1. Claude Code           - Anthropic's Claude AI coding assistant
  2. Aider                 - AI pair programming in your terminal
  3. Continue              - Open-source AI code assistant
  4. GitHub Copilot CLI    - GitHub Copilot command-line interface
  5. Codeium               - Free AI code completion
  6. Cursor                - AI-powered code editor
  7. Qodo                  - AI test generation and code quality
  8. Windsurf              - Codeium's AI-powered IDE
  9. Ollama                - Run large language models locally
 10. LM Studio CLI         - Local LLM management
 11. Sourcegraph Cody      - Sourcegraph's AI coding assistant
 12. Amazon Q              - Amazon's AI developer assistant
 13. Tabnine CLI           - AI code completion
 14. Pieces CLI            - AI-powered code snippet manager
 15. Mentat                - AI coding assistant with context awareness
 16. GPT Engineer          - AI engineer that builds entire codebases
 17. Smol Developer        - Smallest AI developer
 18. Auto-GPT              - Autonomous GPT-4 agent
 19. Open Interpreter      - Natural language interface for computers
 20. Sweep AI              - AI-powered code reviewer
EOF

echo ""
echo "Tools with MCP server integration (for Claude Desktop):"
echo "  • Claude Code (filesystem, github)"
echo "  • Aider (git)"
echo "  • Continue (sqlite)"
echo "  • Cursor (typescript, python)"
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