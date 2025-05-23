package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"mcp-prompt-server/internal/mcp"
	"mcp-prompt-server/internal/prompt"
)

// StdioServer 标准输入输出服务器
type StdioServer struct {
	mcpServer     *mcp.Server
	promptManager *prompt.Manager
	reader        *bufio.Reader
	writer        io.Writer
}

// New 创建新的stdio服务器
func New(mcpServer *mcp.Server, promptManager *prompt.Manager) *StdioServer {
	return &StdioServer{
		mcpServer:     mcpServer,
		promptManager: promptManager,
		reader:        bufio.NewReader(os.Stdin),
		writer:        os.Stdout,
	}
}

// Start 启动服务器
func (s *StdioServer) Start() error {
	log.Println("MCP Prompt Server is running on stdio...")

	for {
		// 读取一行输入
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Received EOF, shutting down...")
				break
			}
			log.Printf("Error reading input: %v", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 处理请求
		if err := s.handleRequest(line); err != nil {
			log.Printf("Error handling request: %v", err)
		}
	}

	return s.promptManager.Close()
}

// handleRequest 处理MCP请求
func (s *StdioServer) handleRequest(requestLine string) error {
	var request mcp.MCPRequest
	if err := json.Unmarshal([]byte(requestLine), &request); err != nil {
		return s.sendError(nil, -32700, "Parse error", err.Error())
	}

	switch request.Method {
	case "initialize":
		return s.handleInitialize(&request)
	case "tools/list":
		return s.handleListTools(&request)
	case "tools/call":
		return s.handleCallTool(&request)
	case "notifications/initialized":
		// 忽略初始化通知
		return nil
	default:
		return s.sendError(request.ID, -32601, "Method not found", fmt.Sprintf("Unknown method: %s", request.Method))
	}
}

// handleInitialize 处理初始化请求
func (s *StdioServer) handleInitialize(request *mcp.MCPRequest) error {
	serverInfo := s.mcpServer.GetServerInfo()

	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": serverInfo,
	}

	return s.sendResponse(request.ID, result)
}

// handleListTools 处理工具列表请求
func (s *StdioServer) handleListTools(request *mcp.MCPRequest) error {
	tools := s.mcpServer.ListTools()
	result := mcp.ListToolsResult{Tools: tools}
	return s.sendResponse(request.ID, result)
}

// handleCallTool 处理工具调用请求
func (s *StdioServer) handleCallTool(request *mcp.MCPRequest) error {
	// 解析参数
	paramsBytes, err := json.Marshal(request.Params)
	if err != nil {
		return s.sendError(request.ID, -32602, "Invalid params", err.Error())
	}

	var params mcp.ToolCallParams
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return s.sendError(request.ID, -32602, "Invalid params", err.Error())
	}

	// 调用工具
	result, err := s.mcpServer.CallTool(params.Name, params.Arguments)
	if err != nil {
		return s.sendError(request.ID, -32603, "Internal error", err.Error())
	}

	return s.sendResponse(request.ID, result)
}

// sendResponse 发送响应
func (s *StdioServer) sendResponse(id interface{}, result interface{}) error {
	response := mcp.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	return s.sendJSON(response)
}

// sendError 发送错误响应
func (s *StdioServer) sendError(id interface{}, code int, message, data string) error {
	response := mcp.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &mcp.MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	return s.sendJSON(response)
}

// sendJSON 发送JSON数据
func (s *StdioServer) sendJSON(data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = fmt.Fprintf(s.writer, "%s\n", string(jsonBytes))
	return err
}
