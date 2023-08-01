PRAGMA foreign_keys = ON;
SELECT load_extension('./crypto');

-- author replaced by hash
-- hashs is actually hash

CREATE TABLE IF NOT EXISTS messages
(
id INTEGER  PRIMARY KEY AUTOINCREMENT,

content TEXT NOT NULL UNIQUE,

hash TEXT NOT NULL UNIQUE
CHECK(hash = sha1(content || COALESCE(parent,''))),

signature TEXT NOT NULL,

parent TEXT,
FOREIGN KEY(parent) REFERENCES messages(hash)
);

CREATE UNIQUE INDEX  IF NOT EXISTS parent_unique ON messages (
       ifnull(parent, '')
);


-- data shite

INSERT INTO messages (content,author, hash) VALUES ('first message','Subodh', sha1('first message'));

WITH parent AS (SELECT hash FROM messages ORDER BY id DESC LIMIT 1)
INSERT INTO messages (content,author, parent, hash) VALUES (
       'second message',
       'Subodh',
       (SELECT hash from parent),
       sha1('second message'||(SELECT hash from parent))
);


WITH parent AS (SELECT hash FROM messages ORDER BY id DESC LIMIT 1)
INSERT INTO messages (content,author, parent, hash) VALUES (
       'third message',
       'Subodh',
       (SELECT hash from parent),
       sha1('third message'||(SELECT hash from parent))
);


WITH parent AS (SELECT hash FROM messages ORDER BY id DESC LIMIT 1)
INSERT INTO messages (content,author, parent, hash) VALUES (
       'fourth message',
       'Subodh',
       (SELECT hash from parent),
       sha1('fourth message'||(SELECT hash from parent))
);

INSERT INTO messages (content,author, hash) VALUES ('first message','NotSubodh', sha1('first message'));

WITH parent AS (SELECT hash FROM messages ORDER BY id DESC LIMIT 1)
INSERT INTO messages (content,author, parent, hash) VALUES (
       'second message',
       'NotSubodh',
       (SELECT hash from parent),
       sha1('second message'||(SELECT hash from parent))
);


WITH parent AS (SELECT hash FROM messages ORDER BY id DESC LIMIT 1)
INSERT INTO messages (content,author, parent, hash) VALUES (
       'third message',
       'NotSubodh',
       (SELECT hash from parent),
       sha1('third message'||(SELECT hash from parent))
);


WITH parent AS (SELECT hash FROM messages ORDER BY id DESC LIMIT 1)
INSERT INTO messages (content,author, parent, hash) VALUES (
       'fourth message',
       'NotSubodh',
       (SELECT hash from parent),
       sha1('fourth message'||(SELECT hash from parent))
);
