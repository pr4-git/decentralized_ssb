package domain

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// MAJOR BEWARE!!!!
// You NEED crypto extension in same dir where you are running from
const loadCryptoExt string = `SELECT load_extension('./crypto');`
const featureFlags string = `PRAGMA foreign_keys = ON;`

// ID is ULID
// Content is our messages's content
// Hash is the sha1 hash of content derived at sqlite side
// (sha1 is generally insecure for passwords but we don't need that here)
// Signature is the ed25519 signature of the author
// Parent is the message preceeding the current one
// (if parent is nil, it means it's the first in thread)
const createMsgTable string = `
CREATE TABLE IF NOT EXISTS messages
(
id TEXT PRIMARY KEY,

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
`

func InitDB(dsn string) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalln("Unable to open database. Cannot recover from this error.")
	}
	db.Exec(loadCryptoExt, nil)
	db.Exec(featureFlags, nil)
	db.Exec(createMsgTable, nil)

}
