package dba

import (
	"crypto/ed25519"
	"log"

	"github.com/jmoiron/sqlx"
)

type Peer struct {
	Pubkey      ed25519.PublicKey `db:"pubkey" json:"pubkey"`
	NetworkAddr string            `db:"networkaddr" json:"address"`
	Alias       string            `db:"alias" json:"alias"`
}

func (peer *Peer) Follow(db *sqlx.DB) error {
	queryStr := `
	INSERT OR REPLACE INTO peers(pubkey,networkaddr, alias)
	VALUES ($1,$2, $3)
	`
	log.Println(peer)
	tx := db.MustBegin()
	_, err := tx.Exec(queryStr, peer.Pubkey, peer.NetworkAddr, peer.Alias)
	if err != nil {
		tx.Rollback()
	} else {
		err = tx.Commit()
	}

	return err
}

func FetchPeerList(db *sqlx.DB) ([]Peer, error) {
	queryStr := `
	SELECT *
	FROM peers
	`

	list := []Peer{}
	err := db.Select(&list, queryStr)
	return list, err
}
