package main

import "github.com/go-chi/chi/v5"

func setupRoutes(r chi.Router, deps *Deps) {

	// all posts that you have made
	r.Get("/@me", deps.ViewWallHandler)

	// view/add peers
	r.Get("/peers", deps.ViewPeerList)
	r.Post("/peers", deps.FollowPeer)

	// view or create posts from peers
	r.Get("/posts", deps.ViewFeedHandler)
	r.Get("/posts/user", deps.ViewAccountPostsHandler)
	r.Post("/posts", deps.PostMessageHandler)

	// Sync all posts with a remote
	r.Post("/sync", deps.RunSyncClient)
}
