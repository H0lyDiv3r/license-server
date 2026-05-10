package main

import (
	"client/internals/state"
	"context"
)

// App struct
type App struct {
	ctx   context.Context
	state *state.AppState
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context, state *state.AppState) {
	a.ctx = ctx
	a.state = state
}
func (a *App) IsLoggedIn() bool {
	return a.state.AuthToken != ""
}
