package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

// Prompt 表示一个prompt模板
type Prompt struct {
	Name        string     `yaml:"name" json:"name"`
	Description string     `yaml:"description" json:"description"`
	Arguments   []Argument `yaml:"arguments" json:"arguments"`
	Messages    []Message  `yaml:"messages" json:"messages"`
}

// Argument 表示prompt的参数
type Argument struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Required    bool   `yaml:"required" json:"required"`
	Type        string `yaml:"type" json:"type"`
}

// Message 表示prompt的消息
type Message struct {
	Role    string  `yaml:"role" json:"role"`
	Content Content `yaml:"content" json:"content"`
}

// Content 表示消息内容
type Content struct {
	Type string `yaml:"type" json:"type"`
	Text string `yaml:"text" json:"text"`
}

// Execute 执行prompt，替换参数并返回最终内容
func (p *Prompt) Execute(args map[string]interface{}) (string, error) {
	var result strings.Builder

	// 处理所有用户消息
	for _, message := range p.Messages {
		if message.Role == "user" && message.Content.Type == "text" {
			content := message.Content.Text

			// 替换参数占位符 {{param}}
			content = p.replaceParameters(content, args)

			result.WriteString(content)
			result.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(result.String()), nil
}

// replaceParameters 替换内容中的参数占位符
func (p *Prompt) replaceParameters(content string, args map[string]interface{}) string {
	// 匹配 {{param}} 格式的占位符
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取参数名
		paramName := strings.Trim(match, "{}")

		// 查找参数值
		if value, exists := args[paramName]; exists {
			return fmt.Sprintf("%v", value)
		}

		// 如果参数不存在，保持原样
		return match
	})
}

// Validate 验证prompt配置的有效性
func (p *Prompt) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("prompt name cannot be empty")
	}

	if len(p.Messages) == 0 {
		return fmt.Errorf("prompt must have at least one message")
	}

	// 检查是否有用户消息
	hasUserMessage := false
	for _, msg := range p.Messages {
		if msg.Role == "user" {
			hasUserMessage = true
			break
		}
	}

	if !hasUserMessage {
		return fmt.Errorf("prompt must have at least one user message")
	}

	return nil
}
