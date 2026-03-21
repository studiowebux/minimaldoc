---
title: MCP Server Documentation
description: Publish static reference pages for MCP servers from JSON manifests
tags:
  - features
  - mcp
  - api
since: "v1.5"
---

# MCP Server Documentation

Generate a static reference page for your MCP server — tools, resources, and prompts — directly from a JSON manifest file. No server required at runtime.

## How It Works

At build time, minimaldoc reads a `.mcp.json` manifest file from your docs directory and generates a self-contained HTML reference page under `/mcp/`. The output is pure static HTML with no JavaScript dependencies.

```
minimaldoc build --mcp
```

Generates:

```
public/
└── mcp/
    └── my-server/
        └── index.html   # Tool list, parameters, resources, prompts
```

## Manifest Format

Create a `.mcp.json` file describing your MCP server's capabilities. The format mirrors the MCP protocol's list responses:

```json
{
  "name": "My MCP Server",
  "description": "Tools for interacting with my service",
  "tools": [
    {
      "name": "create_item",
      "description": "Create a new item in the database",
      "inputSchema": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "description": "Item name"
          },
          "tags": {
            "type": "array",
            "description": "Optional tags",
            "items": { "type": "string" }
          }
        },
        "required": ["name"]
      }
    }
  ],
  "resources": [
    {
      "uri": "file:///data/config.json",
      "name": "Server Config",
      "description": "Current server configuration",
      "mimeType": "application/json"
    }
  ],
  "prompts": [
    {
      "name": "summarize",
      "description": "Summarize recent activity",
      "arguments": [
        {
          "name": "period",
          "description": "Time period to summarize (e.g. '7d', '30d')",
          "required": true
        }
      ]
    }
  ]
}
```

Place the file anywhere under your docs directory (e.g. `mcp/server.mcp.json`). Auto-discovery picks up all `*.mcp.json` files.

## Generating a Manifest

Dump the capabilities of a running MCP server using the MCP inspector:

```bash
npx @modelcontextprotocol/inspector --json > mcp/server.mcp.json
```

Or write it manually — the format is straightforward JSON.

## Configuration

Enable in `config.yaml`:

```yaml
mcp:
  enabled: true
  spec_files:
    - "mcp/*.mcp.json"
  path: "mcp"
```

Or via CLI flag:

```bash
minimaldoc build --mcp
minimaldoc build --mcp --mcp-dir mcp
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | `false` | Enable MCP docs generation |
| `spec_files` | array of string | `[]` | Manifest file paths or globs relative to docs root |
| `path` | string | `"mcp"` | Output URL path |

When `spec_files` is empty, minimaldoc auto-discovers all `*.mcp.json` files under the docs root.

## Output Structure

Each manifest generates a page at `/{path}/{spec-name}/index.html`.

With a single manifest, `/mcp/index.html` redirects directly to the spec page.

With multiple manifests, `/mcp/index.html` lists all available servers.

## What Gets Rendered

The generated page shows:

- **Tools** — each tool with its name, description, and a parameter table (name, type, required/optional, description). Nested object properties are rendered as sub-tables.
- **Resources** — table of all resources with URI, name, MIME type, and description.
- **Prompts** — each prompt with its name, description, and argument table.

Sections with no entries are omitted.
