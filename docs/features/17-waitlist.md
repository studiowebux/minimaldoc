---
title: Waitlist
description: Pre-launch landing page with email signup
tags:
  - features
  - waitlist
  - landing
---

# Waitlist

A pre-launch landing page with email signup form, social links, and privacy policy link.

## Configuration

```yaml
waitlist:
  enabled: true
  title: Coming Soon
  tagline: Sign up for early access
  newsletter_endpoint: https://api.example.com/subscribe
  site_id: my-project
  success_message: "Thanks for signing up! We'll notify you when we launch."
  privacy_url: /legal/privacy.html
  social_links:
    - name: github
      url: https://github.com/org/repo
    - name: discord
      url: https://discord.gg/invite
    - name: twitter
      url: https://twitter.com/handle
```

## Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `false` | Enable the waitlist page |
| `title` | string | `""` | Page headline |
| `tagline` | string | `""` | Subtitle below the headline |
| `newsletter_endpoint` | string | `""` | POST endpoint for the email signup form |
| `site_id` | string | `""` | Site identifier sent as part of the signup payload |
| `success_message` | string | `""` | Message shown after successful signup |
| `privacy_url` | string | `""` | Link to privacy policy (shown below the form) |
| `social_links` | array | `[]` | Social links with `name` and `url` |

## Signup Form

The form submits a POST request to `newsletter_endpoint` with:

```json
{
  "email": "user@example.com",
  "site_id": "my-project"
}
```

On success, the form is replaced with the `success_message`.

## Social Links

Available icon names: `github`, `twitter`, `discord`, `linkedin`, `mastodon`, `youtube`, `rss`, `email`, `website`.

## When to Use

- Product not yet launched but you want to collect interest
- Beta signups before public release
- Landing page while documentation is being written

The waitlist page replaces the normal site layout with a focused single-page design.
