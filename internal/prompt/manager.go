package prompt

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// Manager 管理所有prompt模板
type Manager struct {
	promptsDir string
	prompts    map[string]*Prompt
	mutex      sync.RWMutex
	watcher    *fsnotify.Watcher
}

// NewManager 创建新的prompt管理器
func NewManager(promptsDir string) *Manager {
	return &Manager{
		promptsDir: promptsDir,
		prompts:    make(map[string]*Prompt),
	}
}

// LoadPrompts 加载所有prompt文件
func (m *Manager) LoadPrompts() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 清空现有prompts
	m.prompts = make(map[string]*Prompt)

	// 遍历prompts目录
	err := filepath.WalkDir(m.promptsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 只处理yaml和json文件
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		// 加载prompt文件
		prompt, err := m.loadPromptFile(path)
		if err != nil {
			log.Printf("Warning: Failed to load prompt file %s: %v", path, err)
			return nil // 继续处理其他文件
		}

		// 验证prompt
		if err := prompt.Validate(); err != nil {
			log.Printf("Warning: Invalid prompt in file %s: %v", path, err)
			return nil
		}

		// 检查name冲突
		if _, exists := m.prompts[prompt.Name]; exists {
			log.Printf("Warning: Duplicate prompt name '%s' in file %s, skipping", prompt.Name, path)
			return nil
		}

		// 添加到管理器
		m.prompts[prompt.Name] = prompt

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk prompts directory: %w", err)
	}

	log.Printf("Successfully loaded %d prompts from %s", len(m.prompts), m.promptsDir)

	// 启动文件监控
	if err := m.startWatching(); err != nil {
		log.Printf("Warning: Failed to start file watching: %v", err)
	}

	return nil
}

// loadPromptFile 加载单个prompt文件
func (m *Manager) loadPromptFile(filePath string) (*Prompt, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var prompt Prompt

	// 根据文件扩展名选择解析方式
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".json" {
		if err := json.Unmarshal(data, &prompt); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else {
		// 默认使用YAML解析
		if err := yaml.Unmarshal(data, &prompt); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	}

	return &prompt, nil
}

// GetPrompts 获取所有prompts
func (m *Manager) GetPrompts() []*Prompt {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	prompts := make([]*Prompt, 0, len(m.prompts))
	for _, prompt := range m.prompts {
		prompts = append(prompts, prompt)
	}

	return prompts
}

// GetPrompt 根据名称获取prompt
func (m *Manager) GetPrompt(name string) (*Prompt, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	prompt, exists := m.prompts[name]
	return prompt, exists
}

// GetPromptNames 获取所有prompt名称
func (m *Manager) GetPromptNames() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.prompts))
	for name := range m.prompts {
		names = append(names, name)
	}

	return names
}

// startWatching 启动文件监控，支持热重载
func (m *Manager) startWatching() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}

	m.watcher = watcher

	// 添加prompts目录到监控
	if err := watcher.Add(m.promptsDir); err != nil {
		return fmt.Errorf("failed to watch prompts directory: %w", err)
	}

	// 启动监控goroutine
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 只处理写入和创建事件
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					// 检查是否是支持的文件类型
					ext := strings.ToLower(filepath.Ext(event.Name))
					if ext == ".yaml" || ext == ".yml" || ext == ".json" {
						log.Printf("Detected change in %s, reloading prompts...", event.Name)
						if err := m.LoadPrompts(); err != nil {
							log.Printf("Failed to reload prompts: %v", err)
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("File watcher error: %v", err)
			}
		}
	}()

	log.Printf("Started file watching for %s", m.promptsDir)
	return nil
}

// Close 关闭管理器，清理资源
func (m *Manager) Close() error {
	if m.watcher != nil {
		return m.watcher.Close()
	}
	return nil
}

// Stats 获取统计信息
func (m *Manager) Stats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_prompts":  len(m.prompts),
		"prompts_dir":    m.promptsDir,
		"watching_files": m.watcher != nil,
	}

	// 按类型统计
	argumentCounts := make(map[int]int)
	for _, prompt := range m.prompts {
		argCount := len(prompt.Arguments)
		argumentCounts[argCount]++
	}
	stats["argument_distribution"] = argumentCounts

	return stats
}
