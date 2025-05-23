#!/bin/bash

echo "🎯 get_prompt_names 功能演示"
echo "=============================="
echo

echo "1️⃣ 启动MCP服务器并调用get_prompt_names..."

# 创建测试输入
cat > /tmp/mcp_test_input.json << EOF
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_prompt_names","arguments":{}}}
EOF

echo "2️⃣ 发送请求到MCP服务器..."
echo

# 运行测试并捕获输出
./bin/mcp-prompt-server < /tmp/mcp_test_input.json 2>/tmp/mcp_stderr.log | tail -1 > /tmp/mcp_result.json

echo "3️⃣ 服务器启动日志:"
echo "─────────────────────"
cat /tmp/mcp_stderr.log | grep -E "(Successfully loaded|Registered|Starting)"
echo

echo "4️⃣ get_prompt_names 返回结果:"
echo "─────────────────────────────"

# 解析并美化输出 JSON
if command -v jq &> /dev/null; then
    cat /tmp/mcp_result.json | jq -r '.result.content[0].text'
else
    # 如果没有jq，使用简单的文本处理
    python3 -c "
import json
import sys
try:
    data = json.load(open('/tmp/mcp_result.json'))
    content = data['result']['content'][0]['text']
    print(content)
except:
    print('解析响应时出错')
"
fi

echo "─────────────────────────────"
echo

echo "5️⃣ 功能特点:"
echo "✨ get_prompt_names 工具的优势:"
echo "   🔸 无参数调用: 不需要提供任何参数"
echo "   🔸 实时数据: 返回当前内存中加载的所有prompt"
echo "   🔸 格式清晰: 以列表形式展示，包含总数统计"
echo "   🔸 并发安全: 支持多客户端同时调用"
echo "   🔸 热重载: 反映文件系统的最新变化"

echo
echo "6️⃣ 使用场景:"
echo "   📝 开发阶段: 快速查看可用的prompt工具"
echo "   🔍 工具发现: 帮助用户选择合适的prompt"
echo "   📊 监控统计: 了解prompt库的规模"
echo "   🧪 测试验证: 确认prompt是否正确加载"

# 清理临时文件
rm -f /tmp/mcp_test_input.json /tmp/mcp_result.json /tmp/mcp_stderr.log

echo
echo "🎉 演示完成！" 