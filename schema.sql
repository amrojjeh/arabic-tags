CREATE DATABASE IF NOT EXISTS arabic_tags;

CREATE TABLE IF NOT EXISTS arabic_tags.excerpt (
    id BINARY(16) NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    grammar JSON NOT NULL,
    technical JSON NOT NULL,
    c_locked BOOLEAN NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

