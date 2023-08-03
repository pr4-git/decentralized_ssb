package dba

import (
	"crypto/ed25519"

	"github.com/jmoiron/sqlx"
)

// CREATE TABLE IF NOT EXISTS profiles
// (
// pubkey BLOB PRIMARY KEY,
// privkey BLOB UNIQUE,
// friendlyName TEXT,
// owned BOOL,
// );

type Profile struct {
	PublicKey ed25519.PublicKey `db:"publickey" json:"publickey"`
	Username  string            `db:"Username" json:"username"`
}

func NewProfile(pubkey ed25519.PublicKey, Username string) *Profile {
	return &Profile{
		pubkey,
		Username,
	}
}

// FollowProfile creates database entry of the profiles you follow
func (prof *Profile) FollowProfile(db *sqlx.DB) error {
	queryStr := `
	INSERT INTO following (pubkey, Username)
	VALUES ($1, $2)
	`

	tx := db.MustBegin()
	_, err := tx.Exec(queryStr, prof.PublicKey, prof.Username)
	if err != nil {
		tx.Rollback()
	} else {
		err = tx.Commit()
	}

	return err
}

// Fetch the list of profiles that you follow
func FetchFollowlist(db *sqlx.DB) ([]Profile, error) {
	queryStr := `
	SELECT *
	FROM following
	ORDER BY id ASC
	`
	list := []Profile{}
	err := db.Select(&list, queryStr)
	return list, err
}
