# MinimalDoc AI Documentation Agent

This directory contains reusable AI prompts for generating high-quality documentation in MinimalDoc format.

## Files

- **`minimaldoc-documentation-agent.md`** - Complete agent prompt with templates, best practices, and examples

## Usage

### Option 1: Use with Claude Code (Slash Command)

A slash command version is available at `.claude/commands/minimaldoc-writer.md`:

```bash
/minimaldoc-writer
```

This will activate the documentation writing agent in your current Claude Code session.

### Option 2: Copy as System Prompt

Copy the contents of `minimaldoc-documentation-agent.md` and use it as a system prompt with any AI assistant:

**Claude Desktop/API**:
```json
{
  "system": "<contents of minimaldoc-documentation-agent.md>"
}
```

**ChatGPT Custom Instructions**:
Paste the contents into the "What would you like ChatGPT to know about you" section.

### Option 3: Use as Context

Simply reference the prompt in your conversation:

```
Using the instructions in prompts/minimaldoc-documentation-agent.md, create documentation for [topic]
```

### Option 4: Integrate into CI/CD

Use the prompt with AI API services in your documentation pipeline:

```bash
# Example with Claude API
cat prompts/minimaldoc-documentation-agent.md > agent_prompt.txt
# Use agent_prompt.txt as system prompt for API calls
```

## What It Does

The MinimalDoc Documentation Agent helps you:

1. **Generate Well-Structured Documentation** - Automatically creates proper front matter, headings, and sections
2. **Follow Best Practices** - Ensures clarity, completeness, and consistency
3. **Use Correct Formatting** - Applies MinimalDoc-specific syntax for code blocks, admonitions, and tables
4. **Choose Appropriate Templates** - Selects the right structure for APIs, features, tutorials, or concepts
5. **Optimize for Search** - Includes relevant keywords and tags

## Templates Included

The agent includes templates for:

- **API Documentation** - REST endpoints, parameters, responses
- **Feature Documentation** - Feature guides with usage examples
- **Tutorials** - Step-by-step learning guides
- **Concept Guides** - Explanatory documentation

## Example Usage

**Prompt**:
```
Create documentation for a search feature that supports full-text search with keyboard shortcuts
```

**Agent Output**:
```markdown
---
title: Search Functionality
description: Full-text search with keyboard shortcuts for fast navigation
tags:
  - search
  - keyboard-shortcuts
  - features
---

# Search Functionality

MinimalDoc includes a powerful client-side search feature that allows users to quickly find content across all documentation pages.

[... complete, well-structured documentation ...]
```

## Customization

You can customize the agent by:

1. **Adding Project-Specific Guidelines** - Append your style guide to the prompt
2. **Modifying Templates** - Adjust templates for your use case
3. **Adding Examples** - Include examples from your existing documentation
4. **Changing Tone** - Adjust the writing style (formal/casual)

## Contributing

If you create improved versions or templates, consider:

1. Creating a new template section in the agent prompt
2. Adding examples of high-quality documentation
3. Improving the quality checklist
4. Sharing with the MinimalDoc community

## License

This prompt is part of MinimalDoc and follows the same license.

## Support

For questions or issues:
- Create an issue in the MinimalDoc repository
- Refer to the main documentation
- Join the community discussions
