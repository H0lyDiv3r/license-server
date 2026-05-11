package main

import (
	"client/internals/state"
	"client/internals/store"
	"context"
)

// App struct
type App struct {
	ctx   context.Context
	state *state.AppState
	store *store.Store
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context, state *state.AppState, store *store.Store) {
	a.ctx = ctx
	a.state = state
	a.store = store
}
func (a *App) IsLoggedIn() bool {
	return a.state.AuthToken != ""
}
