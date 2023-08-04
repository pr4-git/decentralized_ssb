package dba

import (
	"database/sql"
	"log"
)

const featureFlags string = `PRAGMA foreign_keys = ON;`

// Feeds are the posts we fetch from others

// ID is ULID
// Content is our feed's content
// Hash is the sha1 hash of content derived at sqlite side
// (sha1 is generally insecure for passwords but we don't need that here)
// Signature is the ed25519 signature of the author
// Parent is the message preceeding the current one
// (if parent is nil, it means it's the first in thread)
const createFeedTable string = `
CREATE TABLE IF NOT EXISTS feed
(
id TEXT PRIMARY KEY,

content TEXT NOT NULL,

signature BLOB NOT NULL,

hash BLOB NOT NULL UNIQUE
CHECK(hash = sha1(content|| signature || COALESCE(parent,''))),

parent BLOB,
author BLOB,

FOREIGN KEY(parent) REFERENCES feed(hash)
FOREIGN KEY(author) REFERENCES peers(pubkey)
);


CREATE UNIQUE INDEX  IF NOT EXISTS parent_unique ON feed (
       ifnull(parent, '')
)
`

// Walls are the posts we make ourselves
const createWallTable string = `
CREATE TABLE IF NOT EXISTS wall
(
id TEXT PRIMARY KEY,

content TEXT NOT NULL,

signature BLOB NOT NULL,

hash BLOB NOT NULL UNIQUE
CHECK(hash = sha1(content|| signature || COALESCE(parent,''))),

parent BLOB,
author BLOB,

FOREIGN KEY(parent) REFERENCES wall(hash)
);


CREATE UNIQUE INDEX  IF NOT EXISTS parent_unique ON wall (
       ifnull(parent, '')
)
`

const createPeersTable string = `
CREATE TABLE IF NOT EXISTS peers
(
	pubkey BLOB PRIMARY KEY,
	networkaddr TEXT NOT NULL,
	alias NOT NULL
)
`

// Initialize the database
// create all the tables and stuff
// WARNING: The order of execution is essential to get right
func InitDB(db *sql.DB) {
	_, err := db.Exec(featureFlags, nil)
	if err != nil {
		db.Close()
		log.Fatalf("Unable to execute. Error %v", err)
	}

	_, err = db.Exec(createPeersTable, nil)
	if err != nil {
		db.Close()
		log.Fatalf("Unable to execute createPeersTable. Error %v", err)
	}
	_, err = db.Exec(createFeedTable, nil)
	if err != nil {
		db.Close()
		log.Fatalf("Unable to execute createFeedTable. Error %v", err)
	}
	_, err = db.Exec(createWallTable, nil)
	if err != nil {
		db.Close()
		log.Fatalf("Unable to execute createWallTable. Error %v", err)
	}
}
