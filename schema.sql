CREATE DATABASE IF NOT EXISTS arabic_tags;

USE arabic_tags;

CREATE TABLE IF NOT EXISTS excerpt (
    id BINARY(16) NOT NULL PRIMARY KEY,
    password_hash CHAR(60) NOT NULL,
    title VARCHAR(100) NOT NULL,
    created DATETIME NOT NULL,
    updated DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS word (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    word VARCHAR(30) NOT NULL,
    excerpt_id BINARY(16) NOT NULL,

    FOREIGN KEY (excerpt_id)
        REFERENCES excerpt(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS manuscript (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    content TEXT NOT NULL,
    locked BOOLEAN NOT NULL,
    excerpt_id BINARY(16) NOT NULL,

    FOREIGN KEY (excerpt_id)
        REFERENCES excerpt(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- For: https://github.com/alexedwards/scs/tree/master/mysqlstore
CREATE TABLE sessions (
	token CHAR(43) PRIMARY KEY,
	data BLOB NOT NULL,
	expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
