---
title: Forum (Optional)
description: Community discussion forum with categories, topics, replies, reputation, and moderation
menu_order: 13
tags:
  - backend
  - forum
  - community
since: 1.4.0
---

# Forum

The optional backend server includes a full-featured community forum for discussions, Q&A, and support.

## Overview

| Feature | Description |
|---------|-------------|
| Categories | Organize discussions by topic with icons, colors, hierarchy |
| Topics | Markdown content with tags, pinning, locking, view counts |
| Replies | Threaded posts with solution marking |
| Reputation | Gamification with points for contributions |
| Moderation | Content flagging, user bans, admin controls |
| Notifications | In-app alerts for replies, likes, solutions |
| Search | Full-text search across topics and posts |

## Configuration

Forum settings via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `FORUM_ENABLED` | `true` | Enable forum feature |
| `FORUM_ALLOW_ANONYMOUS` | `true` | Allow viewing without auth |
| `FORUM_REQUIRE_AUTH` | `true` | Require auth to post |
| `FORUM_MAX_TOPICS_PER_DAY` | `10` | Max topics per user per day |
| `FORUM_MAX_POSTS_PER_DAY` | `50` | Max posts per user per day |
| `FORUM_EDIT_WINDOW` | `1h` | Time window for editing posts |
| `FORUM_MODERATION_MODE` | `none` | `none`, `first_post`, `all` |

### Email Notifications

| Variable | Default | Description |
|----------|---------|-------------|
| `FORUM_EMAIL_ENABLED` | `true` | Enable email notifications |
| `FORUM_EMAIL_DIGEST` | `daily` | `daily`, `weekly`, `none` |
| `FORUM_EMAIL_ON_REPLY` | `true` | Send email on reply |
| `FORUM_EMAIL_ON_MENTION` | `true` | Send email on @mention |

### Reputation System

| Variable | Default | Description |
|----------|---------|-------------|
| `FORUM_REPUTATION_ENABLED` | `true` | Enable reputation points |
| `FORUM_REP_TOPIC_CREATE` | `5` | Points for creating topic |
| `FORUM_REP_POST_CREATE` | `2` | Points for posting reply |
| `FORUM_REP_LIKE_RECEIVED` | `1` | Points when liked |
| `FORUM_REP_SOLUTION_MARKED` | `10` | Points for accepted solution |

## Admin Portal

Access forum administration at `/admin/forum`:

| Page | Description |
|------|-------------|
| Dashboard | Stats overview with recent activity |
| Categories | Create, edit, reorder categories |
| Topics | Moderate topics: pin, lock, close, delete |
| Flags | Review reported content |
| Bans | Manage user bans |
| Tags | Create and manage topic tags |

## Public API

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/forum/categories?site_id=X` | List categories |
| GET | `/api/forum/categories/:slug?site_id=X` | Get category |

### Topics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/forum/topics?site_id=X` | List topics |
| GET | `/api/forum/topics/by-slug/:slug?site_id=X` | Get topic with content |
| GET | `/api/forum/topics/by-slug/:slug/posts?site_id=X` | Get topic replies |
| GET | `/api/forum/search?site_id=X&q=query` | Search topics |

Query parameters for listing:
- `category_id`: Filter by category
- `status`: Filter by status (open, locked, closed)
- `q`: Search query
- `limit`: Results per page (max 100)
- `offset`: Pagination offset

### Tags

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/forum/tags?site_id=X` | List all tags |
| GET | `/api/forum/tags/:slug/topics?site_id=X` | Topics by tag |

### Leaderboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/forum/leaderboard?site_id=X` | Top users by reputation |

## Authenticated API

Requires `Authorization: Bearer <token>` header.

### Topics

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/forum/topics?site_id=X` | Create topic |
| PUT | `/api/forum/topics/:id` | Edit own topic |
| DELETE | `/api/forum/topics/:id` | Delete own topic |

Create topic request:

```json
{
  "title": "Topic title",
  "content": "Markdown content",
  "category_id": "uuid",
  "tags": ["tag-slug-1", "tag-slug-2"]
}
```

### Posts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/forum/topics/by-slug/:slug/posts?site_id=X` | Reply to topic |
| PUT | `/api/forum/posts/:id` | Edit own post |
| DELETE | `/api/forum/posts/:id` | Delete own post |

Reply request:

```json
{
  "content": "Markdown content",
  "parent_id": "optional-parent-post-id"
}
```

### Interactions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/forum/topics/:id/like?site_id=X` | Like/unlike topic |
| POST | `/api/forum/posts/:id/like?site_id=X` | Like/unlike post |
| POST | `/api/forum/topics/:id/bookmark?site_id=X` | Bookmark/unbookmark |
| POST | `/api/forum/flag?site_id=X` | Flag content |

Flag request:

```json
{
  "topic_id": "uuid",
  "reason": "spam|offensive|off-topic|other",
  "description": "Optional details"
}
```

### Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/forum/notifications` | Get notifications |
| PUT | `/api/forum/notifications/:id/read` | Mark as read |
| PUT | `/api/forum/notifications/read-all` | Mark all as read |
| GET | `/api/forum/bookmarks` | Get user bookmarks |

## Admin API

Requires editor role or higher.

### Category Management

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/admin/api/forum/categories` | Editor | Create category |
| PUT | `/admin/api/forum/categories/:id` | Editor | Update category |
| DELETE | `/admin/api/forum/categories/:id` | Editor | Delete category |

### Topic Moderation

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/admin/api/forum/topics/:id/pin` | Editor | Pin/unpin topic |
| POST | `/admin/api/forum/topics/:id/lock` | Editor | Lock topic |
| POST | `/admin/api/forum/topics/:id/close` | Editor | Close topic |
| POST | `/admin/api/forum/topics/:id/open` | Editor | Reopen topic |
| DELETE | `/admin/api/forum/topics/:id` | Editor | Delete topic |
| DELETE | `/admin/api/forum/posts/:id` | Editor | Delete post |
| POST | `/admin/api/forum/posts/:id/solution` | Editor | Mark as solution |

### Flags

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/admin/api/forum/flags` | Editor | List pending flags |
| PUT | `/admin/api/forum/flags/:id` | Editor | Resolve flag |

Resolve flag request:

```json
{
  "status": "resolved|dismissed",
  "resolution_note": "Action taken"
}
```

### Tags

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| POST | `/admin/api/forum/tags` | Editor | Create tag |
| PUT | `/admin/api/forum/tags/:id` | Editor | Update tag |
| DELETE | `/admin/api/forum/tags/:id` | Editor | Delete tag |

### User Bans (Admin Only)

| Method | Endpoint | Role | Description |
|--------|----------|------|-------------|
| GET | `/admin/api/forum/bans` | Admin | List bans |
| POST | `/admin/api/forum/bans` | Admin | Ban user |
| DELETE | `/admin/api/forum/bans/:user_id` | Admin | Unban user |

Ban request:

```json
{
  "user_id": "uuid",
  "reason": "Violation description",
  "expires_at": "2024-12-31T23:59:59Z",
  "is_permanent": false
}
```

## Response Formats

### Topic

```json
{
  "id": "uuid",
  "slug": "topic-title-abc123",
  "title": "Topic Title",
  "content": "Markdown content",
  "content_html": "<p>Rendered HTML</p>",
  "status": "open",
  "is_pinned": false,
  "is_solved": false,
  "view_count": 42,
  "like_count": 5,
  "post_count": 3,
  "category_id": "uuid",
  "category_name": "General",
  "category_slug": "general",
  "author_id": "uuid",
  "author_name": "John Doe",
  "author_avatar": "https://...",
  "tags": ["help", "question"],
  "last_post_at": "2024-01-15T10:30:00Z",
  "created_at": "2024-01-15T09:00:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "is_liked": false,
  "is_bookmarked": false
}
```

### Post

```json
{
  "id": "uuid",
  "topic_id": "uuid",
  "content": "Markdown content",
  "content_html": "<p>Rendered HTML</p>",
  "like_count": 2,
  "is_solution": true,
  "parent_id": "optional-parent-uuid",
  "author_id": "uuid",
  "author_name": "Jane Doe",
  "author_avatar": "https://...",
  "edited_at": "2024-01-15T11:00:00Z",
  "editor_name": "Jane Doe",
  "created_at": "2024-01-15T10:30:00Z",
  "is_liked": false
}
```

### Category

```json
{
  "id": "uuid",
  "slug": "general",
  "name": "General Discussion",
  "description": "General topics",
  "color": "#4a90d9",
  "icon": "chat",
  "position": 1,
  "parent_id": null,
  "is_locked": false,
  "topic_count": 42,
  "post_count": 156
}
```

## User Stats

The forum tracks per-user statistics:

| Stat | Description |
|------|-------------|
| `topic_count` | Topics created |
| `post_count` | Posts/replies created |
| `like_received_count` | Likes received on content |
| `like_given_count` | Likes given to others |
| `solution_count` | Answers marked as solution |
| `reputation` | Total reputation points |

## Notification Types

| Type | Trigger |
|------|---------|
| `reply` | New reply to your topic |
| `like` | Someone liked your content |
| `solution` | Your answer marked as solution |
| `topic_update` | New post in watched topic |

## Topic Status

| Status | Description |
|--------|-------------|
| `open` | Normal, accepting replies |
| `locked` | Visible but no new replies |
| `closed` | Resolved, no new replies |

## Rate Limiting

Built-in rate limiting protects against spam:

- Topics: `FORUM_MAX_TOPICS_PER_DAY` per user
- Posts: `FORUM_MAX_POSTS_PER_DAY` per user
- Banned users cannot create content

## Client Integration

Add forum to your static docs site.

### Enable Forum Feature

In your `config.yaml`:

```yaml
analytics:
  enabled: true
  providers:
    - type: minimaldoc
      enabled: true
      config:
        endpoint: "http://localhost:8090"
        site_id: "your-site-id"
        features: "analytics,forum"
```

Or via script tag:

```html
<script
  src="http://your-server/minimaldoc.js"
  data-endpoint="http://your-server"
  data-site-id="your-site-id"
  data-features="analytics,forum"
  defer>
</script>
```

### Display Forum Components

Add containers to your HTML:

```html
<!-- List forum categories -->
<div data-minimaldoc-forum="categories"></div>

<!-- List latest topics -->
<div data-minimaldoc-forum="latest" data-limit="10"></div>

<!-- List topics by category -->
<div data-minimaldoc-forum="topics" data-category="general" data-limit="20"></div>

<!-- Display single topic with replies -->
<div data-minimaldoc-forum="topic" data-topic="topic-slug-here"></div>

<!-- Forum search -->
<div data-minimaldoc-forum="search"></div>
```

### Programmatic API

The MinimalDoc client provides a JavaScript API for forum interactions:

```javascript
// List categories
MinimalDoc.forum.categories().then(response => {
  console.log(response.categories);
});

// List topics
MinimalDoc.forum.topics({ limit: 10, category: 'general' }).then(response => {
  console.log(response.topics);
});

// Get single topic
MinimalDoc.forum.getTopic('topic-slug').then(response => {
  console.log(response.topic);
});

// Get topic posts/replies
MinimalDoc.forum.getPosts('topic-slug', { limit: 50 }).then(response => {
  console.log(response.posts);
});

// Search forum
MinimalDoc.forum.search('query', 20).then(response => {
  console.log(response.topics);
});

// Get leaderboard
MinimalDoc.forum.leaderboard(10).then(response => {
  console.log(response.leaderboard);
});

// Create topic (requires authentication)
MinimalDoc.forum.createTopic({
  title: 'Topic Title',
  content: 'Markdown content',
  category_id: 'uuid',
  tags: ['tag1', 'tag2']
}).then(() => console.log('Created'));

// Reply to topic (requires authentication)
MinimalDoc.forum.createPost('topic-slug', {
  content: 'Reply content'
}).then(() => console.log('Posted'));

// Like/unlike topic
MinimalDoc.forum.likeTopic('topic-id');

// Bookmark topic
MinimalDoc.forum.bookmark('topic-id');

// Get notifications
MinimalDoc.forum.notifications({ unread: true, limit: 20 });

// Get tags
MinimalDoc.forum.tags();
```

## Forum UI Pages

The forum includes standalone HTML pages with full UI, similar to the blog module. These pages use the landing page layout with header/footer and connect to the API dynamically.

### Available Templates

| Template | URL Pattern | Description |
|----------|-------------|-------------|
| `forum.html` | `/forum/` | Main forum page with categories, recent topics, search |
| `forum-topics.html` | `/forum/topics/` | Browse all topics with filters |
| `forum-category.html` | `/forum/category/:slug/` | Category page with topics |
| `forum-topic.html` | `/forum/topic/:slug/` | Single topic with replies |
| `forum-tag.html` | `/forum/tag/:slug/` | Topics filtered by tag |
| `forum-new.html` | `/forum/new/` | Create new topic form |

### Template Location

Templates are at:
```
internal/assets/themes/common/templates/forum/
```

### Required CSS

The forum uses dedicated styles:
```
internal/assets/themes/common/static/css/forum.css
```

Include in your `config.yaml` or link manually:
```html
<link rel="stylesheet" href="{{.BasePath}}/css/forum.css?v={{.Version}}">
```

### Configuration Variables

Templates expect these variables from the Go template context:

| Variable | Description |
|----------|-------------|
| `{{.BasePath}}` | Base URL path |
| `{{.SiteID}}` | Site identifier for API calls |
| `{{.APIEndpoint}}` | Backend API endpoint URL |
| `{{.Version}}` | Cache-busting version string |

### Linking to Forum

Add forum to your navigation in `config.yaml`:

```yaml
navigation:
  - title: Forum
    url: "/forum/"
    icon: chat
```

Or add a link in your landing page:

```html
<a href="/forum/" class="btn btn-primary">Community Forum</a>
```

### Embedding Components

For embedding forum widgets in other pages, use data attributes:

```html
<!-- List forum categories -->
<div data-minimaldoc-forum="categories"></div>

<!-- List latest topics -->
<div data-minimaldoc-forum="latest" data-limit="10"></div>

<!-- List topics by category -->
<div data-minimaldoc-forum="topics" data-category="general" data-limit="20"></div>

<!-- Display single topic with replies -->
<div data-minimaldoc-forum="topic" data-topic="topic-slug-here"></div>

<!-- Forum search -->
<div data-minimaldoc-forum="search"></div>
```

## Database Schema

The forum uses these tables:

| Table | Description |
|-------|-------------|
| `forum_categories` | Discussion categories |
| `forum_topics` | Topic threads |
| `forum_posts` | Replies/posts |
| `forum_tags` | Topic tags |
| `forum_topic_tags` | Topic-tag relationships |
| `forum_likes` | Likes on topics/posts |
| `forum_bookmarks` | User bookmarks |
| `forum_subscriptions` | Topic watch settings |
| `forum_notifications` | User notifications |
| `forum_flags` | Content reports |
| `forum_bans` | User bans |
| `forum_user_stats` | Per-user statistics |
| `forum_badges` | Achievement badges |
| `forum_user_badges` | Awarded badges |
| `forum_reputation_log` | Reputation history |

Migrations are in:
- `internal/server/store/migrations/sqlite/011_add_forum.*.sql`
- `internal/server/store/migrations/postgres/011_add_forum.*.sql`
