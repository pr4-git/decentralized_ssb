package dba

import (
	"log"

	"github.com/jmoiron/sqlx"
)

const createKnownPeersTable string = `
CREATE TABLE IF NOT EXISTS knownPeers
(
	id TEXT NOT NULL,
	pubkey BLOB PRIMARY KEY,
	alias TEXT DEFAULT "Unknown Peer"
)
`

const createGlobalFeedTable string = `
CREATE TABLE IF NOT EXISTS globalfeed
(
	id TEXT PRIMARY KEY,
	content TEXT NOT NULL,
	signature BLOB NOT NULL,
	hash BLOB NOT NULL,
	parent BLOB,
	author_pk BLOB,

	FOREIGN KEY (parent) REFERENCES globalfeed(hash)
);
`

func Migrate(db *sqlx.DB) {
	tx := db.MustBegin()
	defer tx.Commit()

	_, err := tx.Exec(createKnownPeersTable, nil)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Unable to run migration createKnownPeersTable. Err: %s", err)
	}

	_, err = tx.Exec(createGlobalFeedTable, nil)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Unable to run migration createGLobalFeedTable. Err: %s", err)
	}

}
