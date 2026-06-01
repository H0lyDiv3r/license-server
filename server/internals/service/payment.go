package service

import (
	"context"
	"fmt"
	"license-server/internals/domain"
	"license-server/internals/repository"
	"strconv"

	"github.com/stripe/stripe-go/v81"
	stripesession "github.com/stripe/stripe-go/v81/checkout/session"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return PaymentService{repo: repo}
}

func (p *PaymentService) CreateCheckoutSession(ctx context.Context, amount int64, productName string) (string, error) {
	usr := ctx.Value("user").(domain.UserPayload)
	userId := strconv.FormatFloat(usr.UserId, 'f', -1, 64)
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:3000/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
		Metadata: map[string]string{
			"user_id":    userId,
			"user_email": usr.Email,
			"duration":   "180",
		},
	}

	s, err := stripesession.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return s.URL, nil
}

func (p *PaymentService) SavePaymentSession(ctx context.Context, paymentRequest domain.SavePaymentRequest) (*domain.Payment, error) {
	payment, err := p.repo.StorePayment(ctx, paymentRequest)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
