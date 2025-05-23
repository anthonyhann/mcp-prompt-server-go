#!/bin/bash

echo "🧪 测试Cursor MCP连接状态"
echo "================================"

# 测试1: 验证二进制文件存在
echo "📁 检查MCP服务器二进制文件..."
if [ -f "/Users/hanqiang/www/mcp-prompt-server-go/bin/mcp-prompt-server" ]; then
    echo "✅ MCP服务器二进制文件存在"
else
    echo "❌ MCP服务器二进制文件不存在"
    exit 1
fi

# 测试2: 验证配置文件
echo "📄 检查Cursor配置文件..."
if [ -f "~/.cursor/mcp_config.json" ]; then
    echo "✅ Cursor配置文件存在"
    echo "📋 配置内容："
    cat ~/.cursor/mcp_config.json
else
    echo "❌ Cursor配置文件不存在"
    exit 1
fi

# 测试3: 测试MCP服务器响应
echo "🔄 测试MCP服务器响应..."
cd /Users/hanqiang/www/mcp-prompt-server-go
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/mcp-prompt-server > /tmp/mcp_test.json 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ MCP服务器响应正常"
    # 检查管理工具
    if grep -q "reload_prompts\|get_prompt_names" /tmp/mcp_test.json; then
        echo "✅ 管理工具已注册"
    else
        echo "⚠️  管理工具可能未正确注册"
    fi
else
    echo "❌ MCP服务器响应异常"
    exit 1
fi

echo ""
echo "🎯 在Cursor中使用方法："
echo "1. 确保Cursor已重启"
echo "2. 在聊天中输入：'请使用get_prompt_names工具获取所有可用的prompt名称'"
echo "3. 或者输入：'请使用reload_prompts工具重新加载prompts'"
echo ""
echo "💡 如果仍然无法使用，请尝试："
echo "   - 完全关闭并重启Cursor"
echo "   - 检查Cursor的MCP插件是否启用"
echo "   - 查看Cursor的开发者控制台是否有错误信息"

rm -f /tmp/mcp_test.json 