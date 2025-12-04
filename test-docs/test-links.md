---
title: Link Transformation Test
---

# Link Transformation Test

## Test Cases

### 1. Simple relative link
[Link to index](index.md)

### 2. Relative link with path
[Link to nested](nested/page.md)

### 3. Link with fragment
[Link with anchor](index.md#section)

### 4. Link with number prefix (should be stripped)
[Link to numbered](01-getting-started.md)

### 5. External .md link (should not be transformed)
[External GitHub](https://github.com/user/repo/blob/main/README.md)

### 6. Relative parent directory link
[Link to parent](../example-docs/index.md)

### 7. Multiple nested paths with numbers
[Deep link](../example-docs/getting-started/01-installation.md#prerequisites)
