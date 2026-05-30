package payment

import (
	"client/internals/state"
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (p *Payment) CreateCheckout() string {
	url, err := p.Stripe.CreateCheckout(p.state.AuthToken)
	if err != nil {
		fmt.Println("error getting your url", err)
		return ""
	}
	return url
}

func (p *Payment) OpenCheckoutBrowser(url string) {
	runtime.BrowserOpenURL(p.ctx, url)
}
