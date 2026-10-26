-- Users
CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(255),
    role          VARCHAR(20),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- CREATE INDEX idx_users_email ON users (email);

CREATE TABLE category
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT
);

-- Exhibits
CREATE TABLE exhibits
(
    id               BIGSERIAL PRIMARY KEY,
    category_id      BIGINT       REFERENCES category (id) ON DELETE SET NULL,
--     material         VARCHAR(255) REFERENCES material (name) ON DELETE SET NULL,
    save_by          BIGINT       REFERENCES users (id) ON DELETE SET NULL,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    storage_location VARCHAR(255),
    status           VARCHAR(32)  NOT NULL DEFAULT 'pending_approval',
    length           INT          NOT NULL,
    width            INT          NOT NULL,
    height           INT          NOT NULL,
    creation_date    TIMESTAMPTZ  DEFAULT NOW(),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- CREATE INDEX idx_exhibits_status ON exhibits (status);

-- CREATE TABLE material
-- (
--     name VARCHAR(255) PRIMARY KEY
-- );

CREATE TABLE exhibitions
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    start_date  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    end_date    TIMESTAMPTZ  NOT NULL DEFAULT NOW() + INTERVAL '2 weeks'
);

CREATE TABLE exhibits_exhibitions
(
    exhibition_id BIGINT REFERENCES exhibitions (id) ON DELETE CASCADE,
    exhibit_id BIGINT REFERENCES exhibits (id) ON DELETE CASCADE ,
    PRIMARY KEY (exhibition_id, exhibit_id)
);


-- Trigger: update updated_at
CREATE
OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at
= NOW();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
CREATE TRIGGER exhibits_updated_at
    BEFORE UPDATE
    ON exhibits
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
-- CREATE TRIGGER tasks_updated_at
--     BEFORE UPDATE
--     ON tasks
--     FOR EACH ROW EXECUTE PROCEDURE set_updated_at();