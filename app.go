package main

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"ssb-ng/dba"
	"ssb-ng/gossip"
	"ssb-ng/handshake"

	"github.com/jmoiron/sqlx"
)

type App struct {
	ctx  context.Context
	deps AppDeps
}

type AppDeps struct {
	db      *sqlx.DB
	keypair *handshake.EdKeyPair
	appkey  []byte
}

func (dep AppDeps) GetDB() *sqlx.DB                  { return dep.db }
func (dep AppDeps) GetKeyPair() *handshake.EdKeyPair { return dep.keypair }

// NewApp creates a new App application struct
func NewApp(deps AppDeps) *App {
	return &App{deps: deps}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here

}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) CreateNewPost(content string) error {

	post, err := dba.NewPost(a.deps, content)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	err = dba.InsertPostToDB(a.deps, post)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ViewAllPosts() ([]dba.Post, error) {
	posts, err := dba.FetchAllPosts(a.deps)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (a *App) AddNewPeerProfile(pubkey string, alias string) error {
	peer, err := dba.NewPeer(pubkey, alias)
	if err != nil {
		return err
	}

	err = dba.InsertPeerToDB(a.deps, peer)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ViewAllProfiles() ([]dba.Peer, error) {
	peers, err := dba.FetchKnownPeers(a.deps)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (a *App) ViewOneProfile(pubkeyB64 string) (*dba.Peer, error) {
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyB64)
	if err != nil {
		return nil, err
	}
	peer, err := dba.FetchOnePeer(a.deps, ed25519.PublicKey(pubkey))
	if err != nil {
		return nil, err
	}
	log.Println(peer)
	return peer, nil
}

func (a *App) GetMyPublickey() string {
	kp := GetIdentity()
	pkey := base64.StdEncoding.EncodeToString(kp.Public)
	return pkey
}

func (a *App) IsServer() bool {
	return gossip.IsServer
}

func (a *App) JoinGossipNet(brokerpk string) error {
	// check if we're the server
	if gossip.IsServer {
		return errors.New("this instance is the server")
	}

	log.Printf("Joining gossip network with auth: %s", brokerpk)

	serverPubkey, err := base64.StdEncoding.DecodeString(brokerpk)
	if err != nil {
		return err
	}

	err = gossip.ConnectSyncClient(a.deps, a.deps.appkey, serverPubkey)
	return err
}
