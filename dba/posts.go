package dba

import (
	"crypto/ed25519"
	"errors"
	"itsy/ssb/handshake"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

type Post struct {
	ID        string `db:"id" json:"id"`
	Content   string `db:"content" json:"content"`
	Hash      []byte `db:"hash" json:"hash"`
	Signature []byte `db:"signature" json:"signature"`
	Author    []byte `db:"author" json:"author"`
	Parent    []byte `db:"parent" json:"parent"`
}

func CreatePost(db *sqlx.DB, content string, keypair *handshake.EdKeyPair) error {
	var post = &Post{}
	post.ID = ulid.Make().String()
	post.Content = content
	post.Signature = ed25519.Sign(keypair.Secret, []byte(content))
	post.Author = keypair.Public

	queryStr := `
	WITH parent AS (SELECT * FROM wall ORDER BY id DESC LIMIT 1)
	INSERT INTO wall(id, content,signature,author,parent,hash)
	SELECT
       $1,
       $2,
	   $3,
	   $4,

       CASE WHEN (select count(*) from wall) > 0
       THEN (SELECT hash from parent)
       ELSE NULL END,

       CASE WHEN (select count(*) from wall) > 0
       THEN sha1($2 || $3 || (SELECT hash from parent))
       ELSE sha1($2 || $3)
       END
	`

	tx := db.MustBegin()
	_, err := tx.Exec(queryStr, post.ID, post.Content, post.Signature, post.Author)
	if err != nil {
		tx.Rollback()
	} else {
		err = tx.Commit()
	}

	return err
}

func (msg *Post) SyncToFeed(db *sqlx.DB) error {
	queryStr := `
	WITH parent AS (SELECT * FROM feed ORDER BY id DESC LIMIT 1)
	INSERT INTO feed(id, content,signature,author,parent,hash)
	SELECT
       $1,
       $2,
	   $3,
	   $4,

       CASE WHEN (select count(*) from feed) > 0
       THEN (SELECT hash from parent)
       ELSE NULL END,

       CASE WHEN (select count(*) from feed) > 0
       THEN sha1($2 || $3 || (SELECT hash from parent))
       ELSE sha1($2 || $3)
       END
	`
	// first validate the signature
	if !ed25519.Verify(msg.Author, []byte(msg.Content), msg.Signature) {
		return errors.New("invalid or missing signature of author")
	}
	tx := db.MustBegin()
	_, err := tx.Exec(queryStr, msg.ID, msg.Content, msg.Signature, msg.Author)
	if err != nil {
		tx.Rollback()
	} else {
		err = tx.Commit()
	}

	return err
}

func FetchUserPosts(db *sqlx.DB, author ed25519.PublicKey) ([]Post, error) {
	queryStr := `
	SELECT *
	FROM feed
	WHERE author = $1
	ORDER BY id ASC
	`

	list := []Post{}
	err := db.Select(&list, queryStr, author)
	if err != nil {
		log.Printf("%s", err)
	}
	return list, err
}

func FetchWall(db *sqlx.DB) ([]Post, error) {
	queryStr := `
	SELECT *
	FROM wall
	ORDER BY id ASC
	`

	list := []Post{}
	err := db.Select(&list, queryStr)
	return list, err
}

func FetchFeed(db *sqlx.DB) ([]Post, error) {
	queryStr := `
	SELECT *
	FROM feed
	ORDER BY id DESC
	`

	list := []Post{}
	err := db.Select(&list, queryStr)
	return list, err
}
