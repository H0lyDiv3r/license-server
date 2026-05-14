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
func NewApp(state *state.AppState, store *store.Store) *App {
	return &App{state: state, store: store}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

}
func (a *App) IsLoggedIn() bool {
	return a.state.AuthToken != ""
}
