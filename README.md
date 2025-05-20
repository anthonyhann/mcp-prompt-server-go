# MCP Prompt Server

> 基于 Model Context Protocol (MCP) 的智能提示词服务器，支持一键集成到 Cursor、Windsurf 等编辑器，助力内容创作、代码开发、知识管理和AI自动化。

---

## 主要功能

- 📦 **丰富的Prompt模板**：内置多种高质量Prompt，涵盖代码、写作、产品、知识卡片、网页生成、结构化总结等场景。
- 🛠️ **即插即用的MCP工具**：所有Prompt自动注册为MCP工具，支持参数化调用，适配主流编辑器。
- 🔄 **热加载与管理**：支持一键reload，无需重启即可加载新Prompt。
- 🧩 **极易扩展**：只需添加YAML/JSON文件即可扩展新功能，无需改动主程序。
- 🏷️ **支持多语言与多领域**：适合中英文内容、产品、教育、媒体、AI等多种应用。

---

## 目录结构

```
mcp-prompt-server/
├── package.json
├── src/
│   ├── index.js                # 服务器主入口
│   └── prompts/                # 所有Prompt模板目录
│       ├── gen_summarize.yaml
│       ├── gen_title.yaml
│       ├── gen_html_web_page.yaml
│       ├── gen_3d_webpage_html.yaml
│       ├── gen_bento_grid_html.yaml
│       ├── gen_knowledge_card_html.yaml
│       ├── gen_magazine_card_html.yaml
│       ├── gen_prd_prototype_html.yaml
│       ├── ...                # 更多Prompt模板
│   └── 更多Prompt，需要时拿出来/ # 可选扩展Prompt
└── README.md
```

---

## 快速开始

1. **安装依赖**

   ```bash
   npm install
   ```

2. **启动服务器**

   ```bash
   npm start
   ```

   启动后，MCP Prompt Server会自动加载`src/prompts/`目录下所有Prompt模板，并以MCP工具形式对外提供服务。

---

## 如何使用

### 工具集成

#### Cursor

- 编辑 `~/.cursor/mcp_config.json`，添加如下内容（请将路径替换为你实际的项目路径）：

  ```json
  {
    "servers": [
      {
        "name": "Prompt Server",
        "command": "node",
        "args": [
          "/你的文件实际路径/mcp-prompt-server/src/index.js"
        ],
        "transport": "stdio"
      }
    ]
  }
  ```

- 保存后重启Cursor，即可在工具面板中看到所有Prompt工具。

#### Windsurf

- 编辑 `~/.codeium/windsurf/mcp_config.json`，添加：

  ```json
  {
    "mcpServers": {
      "prompt-server": {
        "command": "node",
        "args": ["/path/to/mcp-prompt-server/src/index.js"],
        "transport": "stdio"
      }
    }
  }
  ```

- 刷新Windsurf设置，Prompt Server即刻生效。

#### Raycast

1. 在 Raycast 搜索 `install server（MCP）`

   ![](https://img.t5t6.com/1747728547294-26c78178-6e42-4e02-a7f3-c9bd9cdbc1fe.png)

2. 给MCP输入一个名字，建议简单点，方便以后@使用，比如叫 `prompt`
3. Command 填写 `node`
4. Argument 填写你的 `index.js` 路径地址

   ![](https://img.t5t6.com/1747728622599-82551d14-937b-4e7c-9429-68d72b7036ce.png)

5. 保存即可，Raycast会自动集成MCP Prompt Server。

##### 注意事项
- 未来新增Prompt，可以复制已有模版让AI参考生成YAML文件。
- **模版中的 `arguments: []` 要么为空，要么参数设置为非必填（false），否则Raycast会报错。**
- 如果报错，可以在Raycast中搜"manage server（MCP）"卸载后重装。
- 每次新增Prompt，都需要卸载重装MCP，暂时没找到更优解。

---

## 如何扩展Prompt

1. **新建YAML或JSON文件**，放入`src/prompts/`目录。
2. **模板格式示例**：

   ```yaml
   name: your_prompt_name
   description: 这个Prompt的用途说明
   arguments: []
   messages:
     - role: user
       content:
         type: text
         text: |
           你的Prompt内容，支持参数占位符{{param}}
   ```

3. **热加载Prompt**  
   - 在编辑器中调用`reload_prompts`工具，或重启服务器即可。

---

## 管理与调试

- `reload_prompts`：热加载所有Prompt模板
- `get_prompt_names`：获取当前所有可用Prompt名称

---

## 高级用法与扩展

- 支持多轮对话Prompt、复杂参数、跨语言内容、数据可视化等高级场景
- 可将`src/prompts/更多Prompt，需要时拿出来/`目录下的模板随时复制到主目录启用

---

## 常见问题

- **Prompt未生效？**  
  检查YAML格式、name字段唯一性，并reload或重启服务。
- **参数不生效？**  
  确认arguments字段正确，调用时传递参数名和值。

---

## 贡献与反馈

- 欢迎提交新Prompt、优化建议或Bug反馈！
- 联系作者：向阳乔木

---

## License

MIT

---

如需进一步定制、批量生成Prompt或企业级集成，欢迎联系作者或提交Issue！
