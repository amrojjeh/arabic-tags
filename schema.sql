CREATE DATABASE IF NOT EXISTS arabic_tags;

USE arabic_tags;

CREATE TABLE IF NOT EXISTS excerpt (
    id BINARY(16) NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    grammar JSON NOT NULL,
    technical JSON NOT NULL,
    c_locked BOOLEAN NOT NULL,
    g_locked BOOLEAN NOT NULL,
    c_share BINARY(16) NOT NULL,
    g_share BINARY(16) NOT NULL,
    t_share BINARY(16) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

CREATE INDEX idx_c_share ON excerpt (c_share);
CREATE INDEX idx_g_share ON excerpt (g_share);
CREATE INDEX idx_t_share ON excerpt (t_share);
