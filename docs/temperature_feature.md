# AI 模型温度设置功能

## 功能概述

SSHAI 现在支持配置 AI 模型的温度参数，用于控制 AI 回答的随机性和创造性。

## 配置方法

在 `config.yaml` 文件的 `api` 部分添加 `temperature` 字段：

```yaml
api:
  base_url: "https://ds.openugc.com/v1"
  api_key: "your-api-key"
  default_model: "deepseek-v3"
  timeout: 600
  temperature: 0.7  # AI模型温度设置，控制回答的随机性 (0.0-2.0)
```

## 温度值说明

| 温度值 | 特性 | 适用场景 |
|--------|------|----------|
| 0.0 | 最确定的回答，输出几乎总是相同 | 需要一致性答案的场景 |
| 0.1-0.3 | 较为保守，适合事实性问答 | 技术文档、代码解释 |
| 0.4-0.7 | 平衡创造性和准确性 | 日常对话、问题解答 |
| 0.8-1.0 | 更有创造性，适合创意写作 | 创意内容、头脑风暴 |
| 1.1-2.0 | 极高随机性，输出可能不太连贯 | 实验性用途 |

## 技术实现

1. **配置结构更新**: 在 `pkg/config/config.go` 中的 `API` 结构体添加了 `Temperature` 字段
2. **请求模型更新**: 在 `pkg/models/models.go` 中的 `ChatRequest` 结构体添加了 `Temperature` 字段
3. **客户端集成**: 在 `pkg/ai/client.go` 中更新了请求构建逻辑，当配置了温度值时会自动包含在 API 请求中

## 使用示例

```yaml
# 保守设置 - 适合技术问答
api:
  temperature: 0.3

# 平衡设置 - 适合日常使用
api:
  temperature: 0.7

# 创意设置 - 适合创作内容
api:
  temperature: 1.0
```

## 注意事项

- 温度值为 0 时，AI 会给出最确定的回答
- 温度值过高（>1.5）可能导致回答不够连贯
- 建议根据具体使用场景调整温度值
- 如果不设置温度值或设置为 0，则不会在 API 请求中包含温度参数