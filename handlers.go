package main

import (
	"context"
	"encoding/base64"
	"errors"
	"itsy/ssb/dba"
	"itsy/ssb/handshake"
	securerpc "itsy/ssb/rpc"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	appkey  string
	db      *sqlx.DB
	keypair handshake.EdKeyPair
	ctx     context.Context
}

func (dep *Deps) startup(ctx context.Context) {
	dep.ctx = ctx
}

func (dep *Deps) RunSyncServer() {
	securerpc.RunServer(dep.db, &dep.keypair, dep.appkey)
}

func (dep *Deps) RunSyncClient(networkAddress string) error {
	if networkAddress == "" {
		return errors.New("invalid input to sync")
	} else {
		securerpc.SyncClient(dep.db, &dep.keypair, networkAddress, dep.appkey)
		return nil
	}
}

func (dep *Deps) FollowPeer(networkAddress string, alias string) error {
	peer := dba.Peer{}
	if networkAddress == "" {
		return errors.New("invalid network address")
	}

	split := strings.Split(networkAddress, "|@")
	split[1] = strings.TrimSuffix(split[1], ".ed25519")
	decoded, err := base64.StdEncoding.DecodeString(split[1])
	peer.Pubkey = decoded

	if err != nil {
		return err
	}
	peer.Follow(dep.db)
	return nil
}

func (dep *Deps) ViewPeerList() ([]dba.Peer, error) {
	peers, err := dba.FetchPeerList(dep.db)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (dep *Deps) PostMessageHandler(content string) error {
	err := dba.CreatePost(dep.db, content, &dep.keypair)
	if err != nil {
		return err
	}
	return nil
}

// wall means our posts
func (dep *Deps) ViewWallHandler() ([]dba.Post, error) {
	posts, err := dba.FetchWall(dep.db)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// feed means other's posts
func (dep *Deps) ViewFeedHandler() ([]dba.Post, error) {
	posts, err := dba.FetchFeed(dep.db)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (dep *Deps) ViewAccountPostsHandler(pubkeystr string) ([]dba.Post, error) {

	pubkey, err := base64.StdEncoding.DecodeString(pubkeystr)
	if err != nil {
		return nil, err
	}

	posts, err := dba.FetchUserPosts(dep.db, pubkey)
	if err != nil {
		return nil, err
	}

	return posts, nil

}
