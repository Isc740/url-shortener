-- +goose Up
CREATE TABLE links (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    target_url TEXT NOT NULL,
    shortened_url TEXT UNIQUE,
    password TEXT,
    status TEXT NOT NULL CHECK (status IN ('active', 'expired' , 'disabled', 'banned')),
    expiration_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_links_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE analytics (
    link_id BIGINT,
    total_clicks INT NOT NULL DEFAULT 0,
    unique_clicks INT NOT NULL DEFAULT 0,
    is_human BOOL DEFAULT TRUE,
    device_type TEXT,
    operating_system TEXT,
    browser TEXT,
    country TEXT,
    referrer TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_analysis_link FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE analytics;
DROP TABLE links;
