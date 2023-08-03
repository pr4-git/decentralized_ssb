package main

import (
	"encoding/base64"
	"encoding/json"
	"itsy/ssb/dba"
	"itsy/ssb/handshake"
	securerpc "itsy/ssb/rpc"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (dep *Deps) RunSyncClient(serverKey string) {
	serverPkey, err := base64.RawStdEncoding.DecodeString(serverKey)
	if err != nil {
		log.Printf("Unable to sync with server. %s", err)
	} else {
		securerpc.SyncClient(dep.db, &dep.keypair, serverPkey, dep.appkey)
	}
}

func (dep *Deps) FollowProfileHandler(w http.ResponseWriter, r *http.Request) {
	b64pubkey := r.Context().Value("pubkey").(string)
	username := r.Context().Value("username").(string)
	pubkey, err := base64.RawStdEncoding.DecodeString(b64pubkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}
	err = dba.NewProfile(pubkey, username).FollowProfile(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	w.WriteHeader(http.StatusOK)
}

func (dep *Deps) ViewFollowingListHandler(w http.ResponseWriter, r *http.Request) {
	profiles, err := dba.FetchFollowlist(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profiles)
}

func (dep *Deps) PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	content := r.Context().Value("content").(string)
	err := dba.CreatePost(dep.db, content, &dep.keypair)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	w.WriteHeader(http.StatusOK)
}

// wall means our posts
func (dep *Deps) ViewWallHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := dba.FetchWall(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// feed means frontpage posts form other users
func (dep *Deps) ViewFeedHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := dba.FetchFeed(dep.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (dep *Deps) ViewAccountPosts(w http.ResponseWriter, r *http.Request) {
	userpubkey := chi.URLParam(r, "pubkey")
	pubkey, err := base64.StdEncoding.DecodeString(userpubkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	posts, err := dba.FetchUserPosts(dep.db, pubkey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]error{"error": err})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
