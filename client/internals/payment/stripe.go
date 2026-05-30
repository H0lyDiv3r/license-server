package payment

import (
	"client/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Stripe struct {
	ctx context.Context
}

func NewStripe(ctx context.Context) Stripe {
	return Stripe{ctx: ctx}
}

func (s *Stripe) CreateCheckout(token string) (string, error) {

	req, err := http.NewRequest("POST", "http://localhost:3000/payment/stripe", nil)

	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var errResp map[string]string
		_ = json.NewDecoder(res.Body).Decode(&errResp)
		if msg := errResp["error"]; msg != "" {
			return "", errors.New(msg)
		}
		return "", fmt.Errorf("payment failed %w", res.Status)
	}

	var result domain.StripeRes
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Println("payment success", result)

	return result.Url, nil
}
