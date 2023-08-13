package dba

import (
	"crypto/ed25519"
	"crypto/sha256"
	"database/sql"
	"log"
	"ssb-ng/handshake"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

var postCount int = 0

type Dependency interface {
	GetDB() *sqlx.DB
	GetKeyPair() *handshake.EdKeyPair
}

type Post struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Hash      []byte `json:"hash"`
	Signature []byte `json:"signature"`
	Author    []byte `db:"author_pk" json:"author"`
	Parent    []byte `json:"parent"`
}

func NewPost(dep Dependency, content string) (*Post, error) {
	keypair := dep.GetKeyPair()
	db := dep.GetDB()
	hasher := sha256.New()

	var parent Post

	err := db.Get(&parent, `
	SELECT *
	FROM globalfeed
	ORDER BY id DESC
	LIMIT 1
	`)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var post = &Post{}
	post.ID = ulid.Make().String()
	post.Content = content
	post.Author = keypair.Public
	post.Signature = ed25519.Sign(keypair.Secret, []byte(content))
	hasher.Write([]byte(post.ID))
	hasher.Write([]byte(post.Content))
	hasher.Write(post.Signature)
	post.Hash = hasher.Sum(nil)
	post.Parent = parent.Hash

	return post, nil
}

func InsertPostToDB(dep Dependency, post *Post) error {
	db := dep.GetDB()

	queryStr := `
	INSERT OR IGNORE INTO globalfeed (id, content, signature, author_pk, parent, hash)
	VALUES($1,$2,$3,$4,$5,$6)
	`
	tx := db.MustBegin()
	defer tx.Commit()
	_, err := tx.Exec(queryStr,
		post.ID,
		post.Content,
		post.Signature,
		post.Author,
		post.Parent,
		post.Hash,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func FetchAllPosts(dep Dependency) ([]Post, error) {
	db := dep.GetDB()

	queryStr := `
	SELECT *
	FROM globalfeed
	ORDER BY ID ASC
	`

	var posts []Post
	err := db.Select(&posts, queryStr)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func FetchPostSeries(dep Dependency, author_pk []byte) ([]Post, error) {
	db := dep.GetDB()

	queryStr := `
	WITH RECURSIVE ancestor(id,hash,parent,content,signature,author_pk) as
	(select globalfeed.id,
		globalfeed.hash,
		globalfeed.parent,
		globalfeed.content,
		globalfeed.signature,
		globalfeed.author_pk
	 from globalfeed
	 where parent is NULL

	 UNION

	select globalfeed.id,
		globalfeed.hash,
		globalfeed.parent,
		globalfeed.content,
		globalfeed.signature,
		globalfeed.author_pk
	 from ancestor,globalfeed

	 where ancestor.hash = feed.parent
	)
	select * 
	from ancestor
	where author_id = $1
	`
	var posts []Post
	err := db.Select(&posts, queryStr)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func RecordsChanged(dep Dependency) (bool, error) {
	db := dep.GetDB()

	queryStr := `
	SELECT count(*)
	FROM globalfeed
	`

	var newCount int
	err := db.Get(&newCount, queryStr)
	if err != nil {
		return false, err
	}

	log.Printf("last count: %d new count: %d", newCount, postCount)

	if newCount != postCount {
		postCount = newCount
		return true, nil
	}

	postCount = newCount
	return false, nil
}
