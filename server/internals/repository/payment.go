package repository

import (
	"context"
	"fmt"
	"license-server/internals/domain"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	GetOne(ctx context.Context, id uint) (*domain.Payment, error)
	StorePayment(ctx context.Context, req domain.SavePaymentRequest) (*domain.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewpaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (p *paymentRepository) GetOne(ctx context.Context, id uint) (*domain.Payment, error) {

	payment, err := gorm.G[domain.Payment](p.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant find payment with id %d: %w", id, err)
	}
	return &payment, err
}

func (p *paymentRepository) StorePayment(ctx context.Context, req domain.SavePaymentRequest) (*domain.Payment, error) {
	// check if license exists
	license, err := gorm.G[domain.License](p.db).Where("id=?", req.LicenseId).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant find license with id %d: %w", req.LicenseId, err)
	}

	payment := &domain.Payment{
		SessionId:     req.SessionId,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Email:         req.Email,
		Status:        req.Status,
		PaymentMethod: req.PaymentMethod,
		MetaData:      req.Metadata,
		License:       license,
	}

	// store payment
	err = gorm.G[domain.Payment](p.db).Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to write key", err.Error())
	}

	return payment, nil
}
