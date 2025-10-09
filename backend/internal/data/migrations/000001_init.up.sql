CREATE TABLE tales (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    summary TEXT,
    parent_id INTEGER NOT NULL DEFAULT 0,
    status_id INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT,
    FOREIGN KEY (parent_id)
    REFERENCES tales (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    FOREIGN KEY (status_id)
    REFERENCES status (id)
        ON UPDATE CASCADE
        ON DELETE SET DEFAULT
);

INSERT INTO tales (
    name,
    summary,
    parent_id,
    status_id,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    "root",
    "This is the root tale",
    0,
    1,
    datetime('now','localtime'),
    datetime('now','localtime'),
    NULL
);

CREATE TABLE chapters (
    id INTEGER PRIMARY KEY,
    content TEXT,
    sentiment REAL,
    tale_id INTEGER NOT NULL,
    FOREIGN KEY (tale_id)
    REFERENCES tale (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE status (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    color TEXT NOT NULL
);

CREATE TABLE tag (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE tale_tag (
    tale_id INTEGER,
    tag_id INTEGER,
    UNIQUE (tale_id, tag_id),
    FOREIGN KEY (tale_id)
    REFERENCES tale (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    FOREIGN KEY (tag_id)
    REFERENCES tag (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE is_similar (
    first_tale_id INTEGER,
    second_tale_id INTEGER,
    UNIQUE (first_tale_id, second_tale_id),
    FOREIGN KEY (first_tale_id)
    REFERENCES tale (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    FOREIGN KEY (second_tale_id)
    REFERENCES tale (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
