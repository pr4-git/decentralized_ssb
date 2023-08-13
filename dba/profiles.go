package dba

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"

	"github.com/oklog/ulid/v2"
)

type Peer struct {
	ID     string            `json:"id"`
	Pubkey ed25519.PublicKey `json:"pubkey"`
	Alias  string            `json:"alias"`
}

func NewPeer(pubkey string, alias string) (*Peer, error) {
	pk, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		return nil, err
	}
	return &Peer{
		ID:     ulid.Make().String(),
		Pubkey: pk,
		Alias:  alias,
	}, nil
}

func InsertPeerToDB(dep Dependency, peer *Peer) error {
	db := dep.GetDB()

	queryStr := `
	INSERT OR IGNORE INTO knownpeers (id, pubkey, alias)
	VALUES($1,$2,$3)
	`
	tx := db.MustBegin()
	defer tx.Commit()
	_, err := tx.Exec(queryStr,
		peer.ID,
		peer.Pubkey,
		peer.Alias,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func FetchKnownPeers(dep Dependency) ([]Peer, error) {
	db := dep.GetDB()

	queryStr := `
	SELECT *
	From knownpeers
	ORDER BY ID ASC
	`

	var peers []Peer
	err := db.Select(&peers, queryStr)
	if err != nil {
		return nil, err
	}

	return peers, nil
}

func FetchOnePeer(dep Dependency, pubkey ed25519.PublicKey) (*Peer, error) {
	db := dep.GetDB()

	queryStr := `
	SELECT *
	From knownpeers
	WHERE pubkey = $1
	LIMIT 1
	`

	var peer Peer
	err := db.Get(&peer, queryStr, pubkey)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &peer, nil
}
