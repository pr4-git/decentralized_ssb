package main

import (
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

type Person struct {
	Username  string `json:"username"`
	PublicKey string `json:"pubkey"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(username string) Person {
	return Person{Username: username, PublicKey: "Unavailable"}
}