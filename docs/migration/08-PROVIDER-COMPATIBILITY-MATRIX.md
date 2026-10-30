# 08 — Provider compatibility matrix

Source: `backend/internal/infra/llm/client.go` registers 16 adapter variants.

- `AdapterOpenAIResponses`
- `AdapterOpenRouterChat`
- `AdapterOpenRouterResponses`
- `AdapterOpenRouterImages`
- `AdapterOpenAIChatCompletions`
- `AdapterOpenAIImageGenerations`
- `AdapterOpenAIImageEdits`
- `AdapterXAIResponses`
- `AdapterXAIImage`
- `AdapterXAIImageEdits`
- `AdapterXAIVideo`
- `AdapterXAIVideoExtensions`
- `AdapterAnthropicMessages`
- `AdapterGoogleGenerateContent`
- `AdapterGoogleImageGeneration`
- `AdapterGeminiInteractions`

Preserve OpenAI Responses and Chat Completions, Anthropic Messages, Google GenerateContent / image generation / Gemini Interactions, xAI Responses/image/edit/video/extensions and OpenRouter Chat/Responses/Images. Keep native request/response shapes, stream events, usage, tool calling, reasoning, prompt cache, multimodal payloads, error classification, retry-after, API key rotation, route priority/weight/fallback and circuit-breaker semantics.

Do not normalize all providers into one OpenAI Chat Completions DTO. Add golden protocol and streaming fixtures for each adapter before G3.
