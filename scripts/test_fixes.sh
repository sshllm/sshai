#!/bin/bash

echo "=== SSH AI 修复功能测试 ==="
echo ""

# 检查服务器是否运行
if ! pgrep -f "./sshai" > /dev/null; then
    echo "启动SSH AI服务器..."
    ./sshai &
    SERVER_PID=$!
    echo "服务器PID: $SERVER_PID"
    sleep 3
else
    echo "SSH AI服务器已在运行"
fi

echo ""
echo "=== 修复内容说明 ==="
echo "1. ✅ 没有匹配模型时显示所有可用模型"
echo "2. ✅ 模型选择界面支持删除输入"
echo "3. ✅ 模型选择界面正确响应回车"
echo "4. ✅ 支持Ctrl+C退出模型选择"
echo ""

echo "=== 测试场景 ==="
echo ""

echo "🔍 场景1: 测试无匹配模型时的处理"
echo "   命令: ssh xyz@localhost -p 2212"
echo "   预期: 显示 '没有找到与用户名 xyz 匹配的模型，显示所有可用模型'"
echo "   然后: 列出所有可用模型供选择"
echo ""

echo "🔍 场景2: 测试模型选择输入功能"
echo "   在模型选择界面："
echo "   - 输入数字 (如: 1, 2, 3)"
echo "   - 使用退格键删除输入"
echo "   - 按回车确认选择"
echo "   - 按Ctrl+C取消选择"
echo ""

echo "🔍 场景3: 测试错误输入处理"
echo "   - 输入无效数字 (如: 0, 99)"
echo "   - 输入非数字字符 (会被忽略)"
echo "   - 空输入后按回车"
echo ""

echo "=== 详细测试步骤 ==="
echo ""

echo "1️⃣ 测试无匹配模型场景:"
echo "   ssh xyz@localhost -p 2212"
echo "   观察是否显示所有模型列表"
echo ""

echo "2️⃣ 测试输入删除功能:"
echo "   在模型选择提示下:"
echo "   - 输入 '12'"
echo "   - 按两次退格键删除"
echo "   - 重新输入 '1'"
echo "   - 按回车确认"
echo ""

echo "3️⃣ 测试错误处理:"
echo "   - 输入 '0' 并按回车 (应提示范围错误)"
echo "   - 输入 'abc1' (只有数字1会被接受)"
echo "   - 直接按回车 (应提示输入有效数字)"
echo ""

echo "4️⃣ 测试中断功能:"
echo "   - 在模型选择时按 Ctrl+C"
echo "   - 应该退出选择并使用默认模型"
echo ""

echo "=== 预期行为验证 ==="
echo ""

echo "✅ 无匹配模型时的消息格式:"
echo "   '没有找到与用户名 'xyz' 匹配的模型，显示所有可用模型:'"
echo ""

echo "✅ 输入处理验证:"
echo "   - 只接受数字字符输入"
echo "   - 退格键正确删除字符"
echo "   - 回车键触发选择验证"
echo "   - 无效选择显示错误提示"
echo ""

echo "✅ 错误消息验证:"
echo "   - '请输入有效的数字: '"
echo "   - '无效输入，请输入数字: '"
echo "   - '请输入 1-N 之间的数字: '"
echo ""

echo "=== 开始测试 ==="
echo "请在新终端中运行以下命令进行测试:"
echo ""

echo "# 测试1: 无匹配模型"
echo "ssh xyz@localhost -p 2212"
echo ""

echo "# 测试2: 有匹配但多个选择"
echo "ssh ai@localhost -p 2212"
echo ""

echo "# 测试3: 无用户名"
echo "ssh localhost -p 2212"
echo ""

echo "在模型选择界面测试以下操作:"
echo "- 输入数字并删除"
echo "- 输入无效数字"
echo "- 按Ctrl+C中断"
echo "- 正常选择模型"
echo ""

echo "按 Ctrl+C 停止服务器"
if [ ! -z "$SERVER_PID" ]; then
    wait $SERVER_PID
fi