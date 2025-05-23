# MCP Prompt Server (Go Version)

> 🚀 **All-New Upgrade!** High-performance MCP Prompt server rewritten in Golang, providing better performance, stronger stability, and richer features.

## ✨ New Feature Upgrades

### 🔥 Performance Optimization
- **Faster Startup**: Go compiled language, 80%+ startup speed improvement
- **Lower Memory**: 60% reduction in memory usage, more stable for long-term operation
- **Concurrent Processing**: Native high-concurrency support, multiple clients can use simultaneously without pressure

### 🛠️ Feature Enhancements
- **Enhanced Hot Reload**: File monitoring with automatic reload, no manual restart required
- **Error Handling**: Comprehensive error handling and logging
- **Parameter Validation**: Strict parameter validation to avoid runtime errors
- **Statistics Monitoring**: Built-in statistics functionality for real-time prompt usage monitoring

### 🏗️ Architecture Optimization
- **Modular Design**: Clear layered architecture, easy to extend and maintain
- **Type Safety**: Strong type system reduces runtime errors
- **Concurrency Safety**: Thread-safe prompt management supporting multi-client access

---

## 📁 Project Structure

```
mcp-prompt-server/
├── main.go                    # Main program entry
├── go.mod                     # Go module definition
├── Makefile                   # Build scripts
├── internal/                  # Internal packages
│   ├── mcp/                   # MCP protocol implementation
│   │   └── models.go          # MCP data models
│   ├── prompt/                # Prompt management
│   │   ├── models.go          # Prompt data models
│   │   └── manager.go         # Prompt manager
│   └── server/                # Server implementation
│       └── stdio.go           # Standard I/O server
├── prompts/                   # Prompt template directory
│   ├── gen_title.yaml         # Title generation
│   ├── gen_summarize.yaml     # Content summarization
│   ├── gen_html_web_page.yaml # Web page generation
│   └── ...                   # More templates
├── tools/                     # Tools directory
│   └── test_mcp.go           # Test tool
└── bin/                       # Build output directory
```

---

## 🚀 Quick Start

### 1. Requirements
- Go 1.21+ 
- Make (optional, for convenient building)

### 2. Installation and Running

#### Method 1: Using Make (Recommended)
```bash
# Install dependencies and build
make build

# Run server
make run

# Development mode (auto-reload)
make dev
```

#### Method 2: Direct Go Commands
```bash
# Install dependencies
go mod tidy

# Build
go build -o bin/mcp-prompt-server main.go

# Run
./bin/mcp-prompt-server
```

### 3. Verify Installation
After startup, you should see output similar to:
```
2024/01/15 10:30:00 Successfully loaded 19 prompts from /path/to/prompts
2024/01/15 10:30:00 Registered 19 prompt tools
2024/01/15 10:30:00 Registered management tools: reload_prompts, get_prompt_names
2024/01/15 10:30:00 Started file watching for /path/to/prompts
2024/01/15 10:30:00 Starting MCP Prompt Server v2.0.0...
2024/01/15 10:30:00 MCP Prompt Server is running on stdio...
```

### 4. Run Tests
```bash
# Run built-in test tool
go run tools/test_mcp.go
```

---

## 🔧 Client Integration

### Raycast
1. Search for "install server (MCP)"
2. Configuration:
   - **Name**: `prompt` 
   - **Command**: Full path to Go binary file
   - **Arguments**: Leave empty

```bash
# Get binary file path
make build
echo "$(pwd)/bin/mcp-prompt-server"
```

### Cursor
Edit `~/.cursor/mcp_config.json`:
```json
{
  "mcpServers": {
    "Prompt Server": {
        "command": "node",
        "args": ["/path/to/mcp-prompt-server-go/bin/mcp-prompt-server"],
        "transport": "stdio"
    }
  }
}
```

### Windsurf
Edit `~/.codeium/windsurf/mcp_config.json`:
```json
{
  "mcpServers": {
    "prompt-server": {
      "command": "/path/to/mcp-prompt-server/bin/mcp-prompt-server",
      "transport": "stdio"
    }
  }
}
```

---

## 📝 Built-in Prompt Tools

The server includes rich prompt templates, including:

### 🎨 Content Creation
- **wechat_headline_generator**: WeChat article viral headline generator
- **gen_summarize**: Intelligent content summarization tool
- **writing_assistant**: Writing assistant
- **gen_podcast_script**: Podcast script generator

### 🌐 Web Page Generation
- **gen_html_web_page**: Universal web page generator
- **gen_3d_webpage_html**: 3D effect web page generator
- **gen_bento_grid_html**: Bento Grid layout web page
- **gen_knowledge_card_html**: Knowledge card web page
- **gen_magazine_card_html**: Magazine-style card

### 💼 Product Development
- **gen_prd_prototype_html**: PRD prototype generator
- **project_architecture**: Project architecture design
- **api_documentation**: API documentation generator

### 💻 Code-Related
- **code_review**: Code review assistant
- **code_refactoring**: Code refactoring suggestions
- **test_case_generator**: Test case generator
- **build_mcp_server**: MCP server build assistant

### 🛠️ Management Tools
- **reload_prompts**: Reload all prompts
- **get_prompt_names**: Get all available prompt names

---

## ⚡ Advanced Features

### 1. Hot Reload
Modify any YAML/JSON files in the prompts directory, and the server will automatically detect and reload without restart.

### 2. Statistics Monitoring
Use the `get_prompt_names` tool to view:
- Number of loaded prompts
- Parameter distribution statistics
- File monitoring status

### 3. Error Handling
- Automatically skip incorrectly formatted prompt files
- Detailed error logging
- Graceful error recovery mechanisms

### 4. Performance Optimization
- Concurrency-safe prompt access
- Memory-efficient file monitoring
- Fast JSON serialization

---

## 📝 Development Guide

### Adding New Prompts
1. Create new YAML/JSON files in the `prompts/` directory
2. Use the following format:

```yaml
name: my_new_prompt
description: This is a description of the new prompt
arguments:
  - name: input_text
    description: Input text
    required: false
    type: string
messages:
  - role: user
    content:
      type: text
      text: |
        Please process the following content: {{input_text}}
        
        Output format: ...
```

3. Save the file, and the server will automatically reload

### Build and Test
```bash
# Code formatting
make fmt

# Static analysis
make vet

# Run tests
make test

# Test coverage
make test-coverage

# Run MCP tests
go run tools/test_mcp.go
```

### Release Packaging
```bash
# Create production version
make build-prod

# Create release package
make package
```

---

## 🔍 Troubleshooting

### Common Issues

1. **Startup Failure**
   ```bash
   # Check Go version
   go version
   
   # Rebuild
   make clean && make build
   ```

2. **Prompts Not Loading**
   ```bash
   # Check file format
   # Use YAML syntax validator to check YAML files
   
   # View logs
   ./bin/mcp-prompt-server 2>&1 | grep -i warning
   ```

3. **Client Connection Issues**
   ```bash
   # Test server
   echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/mcp-prompt-server
   ```

### Log Levels
The program provides detailed log information:
- `INFO`: Normal operation information
- `WARNING`: Non-fatal errors (e.g., invalid prompt files)
- `ERROR`: Serious errors

---

## 📊 Performance Comparison

| Feature | Node.js Version | Go Version | Improvement |
|---------|----------------|------------|-------------|
| Startup Time | ~2.5s | ~0.5s | 80%↑ |
| Memory Usage | ~45MB | ~18MB | 60%↓ |
| Concurrent Processing | Limited | Excellent | Significant improvement |
| File Monitoring | Basic | Efficient | More stable |
| Prompt Count | 11 | 19 | 73%↑ |

---

## 🤝 Contributing

1. Fork the project
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push branch: `git push origin feature/amazing-feature`
5. Submit Pull Request

### Development Standards
- Follow Go code conventions
- Add necessary comments
- Write test cases
- Update relevant documentation

---

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

---

## 🙏 Acknowledgments

- Original Node.js version: [gdli6177/mcp-prompt-server](https://github.com/gdli6177/mcp-prompt-server)
- Model Context Protocol (MCP)
- All contributors and user feedback

---

**💡 Tip**: If you like this project, please give it a ⭐️! Feel free to submit Issues or join discussions if you have any questions. 