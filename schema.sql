CREATE DATABASE IF NOT EXISTS arabic_tags;

USE arabic_tags;

CREATE TABLE IF NOT EXISTS user (
    email VARCHAR(255) NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,

    CONSTRAINT user_username_uc UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS excerpt (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    author_email VARCHAR(255) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,

    FOREIGN KEY (author_email)
        REFERENCES user(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS word (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    word VARCHAR(30) NOT NULL,
    word_pos INTEGER UNSIGNED NOT NULL,
    excerpt_id INTEGER NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,

    FOREIGN KEY (excerpt_id)
        REFERENCES excerpt(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS manuscript (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    content TEXT NOT NULL,
    locked BOOLEAN NOT NULL,
    excerpt_id INTEGER NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL,

    FOREIGN KEY (excerpt_id)
        REFERENCES excerpt(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT manuscript_excerpt_id_uc UNIQUE(excerpt_id)
);

-- For: https://github.com/alexedwards/scs/tree/master/mysqlstore
-- Should be "session" :(
CREATE TABLE sessions (
	token CHAR(43) PRIMARY KEY,
	data BLOB NOT NULL,
	expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
