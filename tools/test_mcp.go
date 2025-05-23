package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("🧪 MCP Prompt Server 测试工具")
	fmt.Println("================================")

	// 启动MCP服务器进程
	cmd := exec.Command("./bin/mcp-prompt-server")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("❌ 无法创建stdin管道: %v\n", err)
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("❌ 无法创建stdout管道: %v\n", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("❌ 无法创建stderr管道: %v\n", err)
		return
	}

	// 启动服务器
	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ 无法启动服务器: %v\n", err)
		return
	}

	// 读取stderr中的启动日志
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Starting MCP Prompt Server") {
				fmt.Printf("✅ 服务器启动成功\n")
			} else if strings.Contains(line, "loaded") && strings.Contains(line, "prompts") {
				fmt.Printf("✅ %s\n", line)
			} else if strings.Contains(line, "ERROR") || strings.Contains(line, "FATAL") {
				fmt.Printf("❌ 错误: %s\n", line)
			}
		}
	}()

	// 给服务器一点时间启动
	time.Sleep(2 * time.Second)

	// 测试用例
	tests := []struct {
		name    string
		request map[string]interface{}
	}{
		{
			name: "初始化请求",
			request: map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"method":  "initialize",
				"params": map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities":    map[string]interface{}{},
					"clientInfo": map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
			},
		},
		{
			name: "获取工具列表",
			request: map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      2,
				"method":  "tools/list",
			},
		},
		{
			name: "获取prompt名称",
			request: map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      3,
				"method":  "tools/call",
				"params": map[string]interface{}{
					"name":      "get_prompt_names",
					"arguments": map[string]interface{}{},
				},
			},
		},
	}

	scanner := bufio.NewScanner(stdout)

	// 执行测试
	for i, test := range tests {
		fmt.Printf("\n🔄 测试 %d: %s\n", i+1, test.name)

		// 发送请求
		requestBytes, _ := json.Marshal(test.request)
		requestStr := string(requestBytes) + "\n"

		if _, err := stdin.Write([]byte(requestStr)); err != nil {
			fmt.Printf("❌ 发送请求失败: %v\n", err)
			continue
		}

		// 读取响应
		if scanner.Scan() {
			response := scanner.Text()

			var result map[string]interface{}
			if err := json.Unmarshal([]byte(response), &result); err != nil {
				fmt.Printf("❌ 解析响应失败: %v\n", err)
				fmt.Printf("原始响应: %s\n", response)
				continue
			}

			if errField, exists := result["error"]; exists {
				fmt.Printf("❌ 服务器返回错误: %v\n", errField)
			} else {
				fmt.Printf("✅ 测试通过\n")
				if test.name == "获取工具列表" {
					if resultField, exists := result["result"]; exists {
						if tools, ok := resultField.(map[string]interface{})["tools"]; ok {
							if toolsArray, ok := tools.([]interface{}); ok {
								fmt.Printf("📝 发现 %d 个工具\n", len(toolsArray))
								for _, tool := range toolsArray {
									if toolMap, ok := tool.(map[string]interface{}); ok {
										if name, exists := toolMap["name"]; exists {
											fmt.Printf("   - %s\n", name)
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			fmt.Printf("❌ 未收到响应\n")
		}

		time.Sleep(500 * time.Millisecond)
	}

	// 关闭连接并终止进程
	stdin.Close()
	cmd.Process.Kill()
	cmd.Wait()

	fmt.Println("\n🎉 测试完成！")
}
