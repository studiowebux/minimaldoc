---
title: Claude Assist
description: Add an AI assistant button to documentation pages
tags:
  - features
  - ai
  - claude
---

# Claude Assist

Add an "Ask Claude" button to each page that opens Claude with the page content as context.

## Configuration

```yaml
claude_assist:
  enabled: true
  label: Ask Claude
  prompt: "Answer based on this documentation page:"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Show the Ask Claude button on each page |
| `label` | string | `"Ask Claude"` | Button text |
| `prompt` | string | `""` | Custom prompt prefix prepended to the page content |

## How It Works

When clicked, the button:

1. Extracts the current page's markdown content (from the companion `.html.md` file)
2. Prepends the optional custom prompt
3. Opens Claude (claude.ai) with the content as context

This requires `enable_llms: true` since the button depends on the generated `.html.md` companion files.

## Use Cases

- Let users ask follow-up questions about complex configuration
- Provide an AI-powered support layer on top of static docs
- Help users translate documentation concepts to their specific use case

## Button Placement

The button appears in the page title row alongside other action buttons (Export PDF, Copy MD). On mobile, buttons wrap below the title.

## Custom Prompts

The `prompt` field lets you customize how Claude interprets the page:

```yaml
# Technical docs
claude_assist:
  enabled: true
  prompt: "You are a technical support engineer. Answer based on this documentation:"

# API reference
claude_assist:
  enabled: true
  label: "Explain this API"
  prompt: "Explain this API endpoint in simple terms:"
```

When no prompt is set, the page content is sent as-is.
