package payment

import (
	"client/internals/state"
	"context"
)

type Payment struct {
	ctx    context.Context
	state  *state.AppState
	Stripe Stripe
}

func NewPayment(s *state.AppState) *Payment {
	return &Payment{
		state: s,
	}
}

func (p *Payment) Startup(ctx context.Context) {
	p.ctx = ctx
	p.Stripe = NewStripe(ctx)
}

func (p *Payment) CreateCheckout() {
	p.Stripe.CreateCheckout(p.state.AuthToken)
}
