package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mcp-prompt-server/internal/mcp"
	"mcp-prompt-server/internal/prompt"
	"mcp-prompt-server/internal/server"
)

// 构建时注入的版本信息
var (
	Version    = "2.0.0"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

const (
	serverName = "mcp-prompt-server"
	promptsDir = "prompts"
)

func main() {
	// 初始化日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// 设置prompts目录路径
	promptsDirPath := filepath.Join(workDir, promptsDir)

	// 确保prompts目录存在
	if err := os.MkdirAll(promptsDirPath, 0755); err != nil {
		log.Fatalf("Failed to create prompts directory: %v", err)
	}

	// 创建prompt管理器
	promptManager := prompt.NewManager(promptsDirPath)

	// 加载所有prompts
	if err := promptManager.LoadPrompts(); err != nil {
		log.Fatalf("Failed to load prompts: %v", err)
	}

	// 创建MCP服务器
	mcpServer := mcp.NewServer(serverName, Version)

	// 注册prompt工具
	if err := registerPromptTools(mcpServer, promptManager); err != nil {
		log.Fatalf("Failed to register prompt tools: %v", err)
	}

	// 注册管理工具
	registerManagementTools(mcpServer, promptManager)

	// 创建服务器实例
	srv := server.New(mcpServer, promptManager)

	// 启动服务器
	log.Printf("Starting MCP Prompt Server v%s (built: %s, commit: %s)...", Version, BuildTime, CommitHash)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// registerPromptTools 注册所有prompt工具
func registerPromptTools(mcpServer *mcp.Server, promptManager *prompt.Manager) error {
	prompts := promptManager.GetPrompts()

	for _, p := range prompts {
		tool := &mcp.Tool{
			Name:        p.Name,
			Description: p.Description,
			Arguments:   buildArgumentSchema(p.Arguments),
			Handler: func(prompt *prompt.Prompt) mcp.ToolHandler {
				return func(args map[string]interface{}) (*mcp.ToolResult, error) {
					return executePrompt(prompt, args)
				}
			}(p),
		}

		if err := mcpServer.RegisterTool(tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %w", p.Name, err)
		}
	}

	log.Printf("Registered %d prompt tools", len(prompts))
	return nil
}

// registerManagementTools 注册管理工具
func registerManagementTools(mcpServer *mcp.Server, promptManager *prompt.Manager) {
	// 重新加载prompts工具
	reloadTool := &mcp.Tool{
		Name:        "reload_prompts",
		Description: "重新加载所有预设的prompts",
		Arguments:   map[string]interface{}{},
		Handler: func(args map[string]interface{}) (*mcp.ToolResult, error) {
			if err := promptManager.LoadPrompts(); err != nil {
				return &mcp.ToolResult{
					Content: []mcp.Content{{
						Type: "text",
						Text: fmt.Sprintf("重新加载失败: %v", err),
					}},
				}, nil
			}

			count := len(promptManager.GetPrompts())
			return &mcp.ToolResult{
				Content: []mcp.Content{{
					Type: "text",
					Text: fmt.Sprintf("成功重新加载了 %d 个prompts。", count),
				}},
			}, nil
		},
	}

	// 获取prompt名称列表工具
	listTool := &mcp.Tool{
		Name:        "get_prompt_names",
		Description: "获取所有可用的prompt名称",
		Arguments:   map[string]interface{}{},
		Handler: func(args map[string]interface{}) (*mcp.ToolResult, error) {
			prompts := promptManager.GetPrompts()
			names := make([]string, len(prompts))
			for i, p := range prompts {
				names[i] = p.Name
			}

			result := fmt.Sprintf("可用的prompts (%d):\n", len(names))
			for _, name := range names {
				result += "- " + name + "\n"
			}

			return &mcp.ToolResult{
				Content: []mcp.Content{{
					Type: "text",
					Text: result,
				}},
			}, nil
		},
	}

	mcpServer.RegisterTool(reloadTool)
	mcpServer.RegisterTool(listTool)

	log.Println("Registered management tools: reload_prompts, get_prompt_names")
}

// buildArgumentSchema 构建参数schema
func buildArgumentSchema(args []prompt.Argument) map[string]interface{} {
	schema := make(map[string]interface{})

	for _, arg := range args {
		argSchema := map[string]interface{}{
			"type":        "string",
			"description": arg.Description,
		}

		if arg.Required {
			argSchema["required"] = true
		}

		schema[arg.Name] = argSchema
	}

	return schema
}

// executePrompt 执行prompt并返回结果
func executePrompt(p *prompt.Prompt, args map[string]interface{}) (*mcp.ToolResult, error) {
	content, err := p.Execute(args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute prompt: %w", err)
	}

	return &mcp.ToolResult{
		Content: []mcp.Content{{
			Type: "text",
			Text: content,
		}},
	}, nil
}
