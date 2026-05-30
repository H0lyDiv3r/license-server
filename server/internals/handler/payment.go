package handler

import (
	"encoding/json"
	"license-server/internals/domain"
	"license-server/internals/service"
	"net/http"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) PaymentHandler {
	return PaymentHandler{
		service: service,
	}
}

func (p *PaymentHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	checkoutURL, err := p.service.CreateCheckoutSession(2000, "fake_product")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domain.CreateCheckoutResponse{Url: checkoutURL})
}
