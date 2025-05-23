package mcp

import (
	"fmt"
)

// Server MCP服务器
type Server struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	tools   map[string]*Tool
}

// Tool MCP工具定义
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Arguments   map[string]interface{} `json:"arguments"`
	Handler     ToolHandler            `json:"-"`
}

// ToolHandler 工具处理函数
type ToolHandler func(args map[string]interface{}) (*ToolResult, error)

// ToolResult 工具执行结果
type ToolResult struct {
	Content []Content `json:"content"`
}

// Content 内容结构
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MCPRequest MCP请求
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse MCP响应
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError MCP错误
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToolCallParams 工具调用参数
type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ListToolsResult 工具列表结果
type ListToolsResult struct {
	Tools []ToolInfo `json:"tools"`
}

// ToolInfo 工具信息
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewServer 创建新的MCP服务器
func NewServer(name, version string) *Server {
	return &Server{
		Name:    name,
		Version: version,
		tools:   make(map[string]*Tool),
	}
}

// RegisterTool 注册工具
func (s *Server) RegisterTool(tool *Tool) error {
	if tool.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if tool.Handler == nil {
		return fmt.Errorf("tool handler cannot be nil")
	}

	s.tools[tool.Name] = tool
	return nil
}

// GetTool 获取工具
func (s *Server) GetTool(name string) (*Tool, bool) {
	tool, exists := s.tools[name]
	return tool, exists
}

// ListTools 列出所有工具
func (s *Server) ListTools() []ToolInfo {
	tools := make([]ToolInfo, 0, len(s.tools))

	for _, tool := range s.tools {
		toolInfo := ToolInfo{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: buildInputSchema(tool.Arguments),
		}
		tools = append(tools, toolInfo)
	}

	return tools
}

// CallTool 调用工具
func (s *Server) CallTool(name string, args map[string]interface{}) (*ToolResult, error) {
	tool, exists := s.GetTool(name)
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	return tool.Handler(args)
}

// GetServerInfo 获取服务器信息
func (s *Server) GetServerInfo() ServerInfo {
	return ServerInfo{
		Name:    s.Name,
		Version: s.Version,
	}
}

// buildInputSchema 构建输入schema
func buildInputSchema(arguments map[string]interface{}) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": arguments,
	}

	// 提取必需的参数
	required := make([]string, 0)
	for name, arg := range arguments {
		if argMap, ok := arg.(map[string]interface{}); ok {
			if req, exists := argMap["required"]; exists {
				if isRequired, ok := req.(bool); ok && isRequired {
					required = append(required, name)
				}
			}
		}
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}
