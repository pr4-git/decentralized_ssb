package main

import (
	"encoding/base64"
	"encoding/json"
	"itsy/ssb/dba"
	"itsy/ssb/handshake"
	securerpc "itsy/ssb/rpc"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	appkey  string
	db      *sqlx.DB
	keypair handshake.EdKeyPair
}

func (dep *Deps) RunSyncServer() {
	securerpc.RunServer(dep.db, &dep.keypair, dep.appkey)
}

func (dep *Deps) RunSyncClient(w http.ResponseWriter, r *http.Request) {
	var form struct {
		NetworkAddr string `json:"address"`
	}

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil || form.NetworkAddr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		securerpc.SyncClient(dep.db, &dep.keypair, form.NetworkAddr, dep.appkey)
	}
}

func (dep *Deps) FollowPeer(w http.ResponseWriter, r *http.Request) {
	var peer dba.Peer

	err := json.NewDecoder(r.Body).Decode(&peer)
	if err != nil || peer.NetworkAddr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}
	log.Printf("Peer: %v", peer)
	split := strings.Split(peer.NetworkAddr, "|@")
	split[1] = strings.TrimSuffix(split[1], ".ed25519")
	peer.Pubkey, err = base64.StdEncoding.DecodeString(split[1])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		peer.Follow(dep.db)
		w.WriteHeader(http.StatusOK)
	}
}

func (dep *Deps) ViewPeerList(w http.ResponseWriter, r *http.Request) {
	peers, err := dba.FetchPeerList(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(peers)
}

func (dep *Deps) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	var form struct{ Content string }

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	log.Printf("Content: %s", form.Content)

	err = dba.CreatePost(dep.db, form.Content, &dep.keypair)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(http.StatusOK)
}

// wall means our posts
func (dep *Deps) ViewWallHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := dba.FetchWall(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// feed means frontpage posts form other users
func (dep *Deps) ViewFeedHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := dba.FetchFeed(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (dep *Deps) ViewAccountPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("triggered!")
	var form struct{ Pubkey string }

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil || form.Pubkey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	pubkey, err := base64.StdEncoding.DecodeString(form.Pubkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	posts, err := dba.FetchUserPosts(dep.db, pubkey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
