-- Forum categories table (hierarchical with one level nesting)
CREATE TABLE forum_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES forum_categories(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    color TEXT,
    icon TEXT,
    position INTEGER DEFAULT 0,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, slug)
);

-- Forum topics table (discussion threads)
CREATE TABLE forum_topics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    category_id UUID REFERENCES forum_categories(id) ON DELETE CASCADE,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    slug TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    status TEXT DEFAULT 'open' CHECK (status IN ('open', 'closed', 'locked', 'archived')),
    is_pinned BOOLEAN DEFAULT FALSE,
    is_solved BOOLEAN DEFAULT FALSE,
    solution_post_id UUID,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    last_post_at TIMESTAMPTZ,
    last_post_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, slug)
);

-- Forum posts table (replies with threading)
CREATE TABLE forum_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    like_count INTEGER DEFAULT 0,
    is_solution BOOLEAN DEFAULT FALSE,
    edited_at TIMESTAMPTZ,
    edited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Forum tags table
CREATE TABLE forum_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    color TEXT,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, slug)
);

-- Forum topic tags (many-to-many)
CREATE TABLE forum_topic_tags (
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES forum_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (topic_id, tag_id)
);

-- Forum likes (topics and posts)
CREATE TABLE forum_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    post_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT forum_likes_target CHECK (
        (topic_id IS NOT NULL AND post_id IS NULL) OR
        (topic_id IS NULL AND post_id IS NOT NULL)
    ),
    UNIQUE(user_id, topic_id),
    UNIQUE(user_id, post_id)
);

-- Forum bookmarks (saved topics)
CREATE TABLE forum_bookmarks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, topic_id)
);

-- Forum subscriptions (watch/mute topics)
CREATE TABLE forum_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    category_id UUID REFERENCES forum_categories(id) ON DELETE CASCADE,
    level TEXT DEFAULT 'watching' CHECK (level IN ('watching', 'tracking', 'muted')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT forum_subscriptions_target CHECK (
        (topic_id IS NOT NULL AND category_id IS NULL) OR
        (topic_id IS NULL AND category_id IS NOT NULL)
    )
);

-- Forum mentions (@mentions tracking)
CREATE TABLE forum_mentions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    mentioned_by UUID REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    post_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Forum notifications (in-app notifications)
CREATE TABLE forum_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('reply', 'mention', 'like', 'solution', 'badge', 'topic_update')),
    title TEXT NOT NULL,
    message TEXT,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    post_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Forum flags (content reports)
CREATE TABLE forum_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    reporter_id UUID REFERENCES users(id) ON DELETE SET NULL,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE CASCADE,
    post_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE,
    reason TEXT NOT NULL CHECK (reason IN ('spam', 'inappropriate', 'harassment', 'off_topic', 'misinformation', 'other')),
    description TEXT,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'resolved', 'rejected')),
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMPTZ,
    resolution_note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT forum_flags_target CHECK (
        (topic_id IS NOT NULL AND post_id IS NULL) OR
        (topic_id IS NULL AND post_id IS NOT NULL)
    )
);

-- Forum bans (user bans with expiration)
CREATE TABLE forum_bans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    banned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reason TEXT NOT NULL,
    expires_at TIMESTAMPTZ,
    is_permanent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, user_id)
);

-- Forum user stats (reputation, counts)
CREATE TABLE forum_user_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reputation INTEGER DEFAULT 0,
    topic_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    like_received_count INTEGER DEFAULT 0,
    like_given_count INTEGER DEFAULT 0,
    solution_count INTEGER DEFAULT 0,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, user_id)
);

-- Forum badges (badge definitions)
CREATE TABLE forum_badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    color TEXT,
    criteria TEXT,
    tier TEXT DEFAULT 'bronze' CHECK (tier IN ('bronze', 'silver', 'gold')),
    is_manual BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(site_id, slug)
);

-- Forum user badges (badges awarded to users)
CREATE TABLE forum_user_badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    badge_id UUID REFERENCES forum_badges(id) ON DELETE CASCADE,
    awarded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    awarded_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, badge_id)
);

-- Forum reputation log (point history for auditing)
CREATE TABLE forum_reputation_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID REFERENCES sites(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    action TEXT NOT NULL CHECK (action IN ('topic_create', 'post_create', 'like_received', 'like_given', 'solution_marked', 'badge_awarded', 'manual_adjust')),
    points INTEGER NOT NULL,
    topic_id UUID REFERENCES forum_topics(id) ON DELETE SET NULL,
    post_id UUID REFERENCES forum_posts(id) ON DELETE SET NULL,
    badge_id UUID REFERENCES forum_badges(id) ON DELETE SET NULL,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add foreign key for solution_post_id after forum_posts is created
ALTER TABLE forum_topics ADD CONSTRAINT fk_solution_post
    FOREIGN KEY (solution_post_id) REFERENCES forum_posts(id) ON DELETE SET NULL;

-- Indexes for forum_categories
CREATE INDEX idx_forum_categories_site ON forum_categories(site_id);
CREATE INDEX idx_forum_categories_parent ON forum_categories(parent_id);
CREATE INDEX idx_forum_categories_position ON forum_categories(site_id, position);

-- Indexes for forum_topics
CREATE INDEX idx_forum_topics_site ON forum_topics(site_id);
CREATE INDEX idx_forum_topics_category ON forum_topics(category_id);
CREATE INDEX idx_forum_topics_author ON forum_topics(author_id);
CREATE INDEX idx_forum_topics_site_status ON forum_topics(site_id, status);
CREATE INDEX idx_forum_topics_pinned ON forum_topics(site_id, is_pinned DESC, created_at DESC);
CREATE INDEX idx_forum_topics_latest ON forum_topics(site_id, last_post_at DESC);
CREATE INDEX idx_forum_topics_slug ON forum_topics(site_id, slug);

-- Indexes for forum_posts
CREATE INDEX idx_forum_posts_topic ON forum_posts(topic_id);
CREATE INDEX idx_forum_posts_parent ON forum_posts(parent_id);
CREATE INDEX idx_forum_posts_author ON forum_posts(author_id);
CREATE INDEX idx_forum_posts_created ON forum_posts(topic_id, created_at);

-- Indexes for forum_tags
CREATE INDEX idx_forum_tags_site ON forum_tags(site_id);
CREATE INDEX idx_forum_tags_usage ON forum_tags(site_id, usage_count DESC);

-- Indexes for forum_topic_tags
CREATE INDEX idx_forum_topic_tags_tag ON forum_topic_tags(tag_id);

-- Indexes for forum_likes
CREATE INDEX idx_forum_likes_user ON forum_likes(user_id);
CREATE INDEX idx_forum_likes_topic ON forum_likes(topic_id);
CREATE INDEX idx_forum_likes_post ON forum_likes(post_id);

-- Indexes for forum_bookmarks
CREATE INDEX idx_forum_bookmarks_user ON forum_bookmarks(user_id);
CREATE INDEX idx_forum_bookmarks_topic ON forum_bookmarks(topic_id);

-- Indexes for forum_subscriptions
CREATE INDEX idx_forum_subscriptions_user ON forum_subscriptions(user_id);
CREATE INDEX idx_forum_subscriptions_topic ON forum_subscriptions(topic_id);
CREATE INDEX idx_forum_subscriptions_category ON forum_subscriptions(category_id);

-- Indexes for forum_mentions
CREATE INDEX idx_forum_mentions_user ON forum_mentions(user_id);
CREATE INDEX idx_forum_mentions_unread ON forum_mentions(user_id, is_read);

-- Indexes for forum_notifications
CREATE INDEX idx_forum_notifications_user ON forum_notifications(user_id);
CREATE INDEX idx_forum_notifications_unread ON forum_notifications(user_id, is_read);
CREATE INDEX idx_forum_notifications_created ON forum_notifications(user_id, created_at DESC);

-- Indexes for forum_flags
CREATE INDEX idx_forum_flags_site ON forum_flags(site_id);
CREATE INDEX idx_forum_flags_status ON forum_flags(site_id, status);
CREATE INDEX idx_forum_flags_pending ON forum_flags(site_id, status) WHERE status = 'pending';

-- Indexes for forum_bans
CREATE INDEX idx_forum_bans_user ON forum_bans(user_id);
CREATE INDEX idx_forum_bans_active ON forum_bans(site_id, user_id) WHERE is_permanent = TRUE OR expires_at > NOW();

-- Indexes for forum_user_stats
CREATE INDEX idx_forum_user_stats_user ON forum_user_stats(user_id);
CREATE INDEX idx_forum_user_stats_reputation ON forum_user_stats(site_id, reputation DESC);

-- Indexes for forum_badges
CREATE INDEX idx_forum_badges_site ON forum_badges(site_id);

-- Indexes for forum_user_badges
CREATE INDEX idx_forum_user_badges_user ON forum_user_badges(user_id);
CREATE INDEX idx_forum_user_badges_badge ON forum_user_badges(badge_id);

-- Indexes for forum_reputation_log
CREATE INDEX idx_forum_reputation_log_user ON forum_reputation_log(user_id);
CREATE INDEX idx_forum_reputation_log_created ON forum_reputation_log(user_id, created_at DESC);
