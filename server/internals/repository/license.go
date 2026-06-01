package repository

import (
	"context"
	"fmt"
	"license-server/internals/domain"

	"gorm.io/gorm"
)

type LicenseRepository interface {
	StoreKey(ctx context.Context, req domain.StoreKeyReq) (domain.License, error)
	GetLicenseById(ctx context.Context, id string) (domain.License, error)
	GetLicenseByKey(ctx context.Context, key string) (*domain.License, error)
	GetLicenseBySessionId(ctx context.Context, sessionId string) (*domain.License, error)
	UpdateLicense(ctx context.Context, license domain.License) (domain.License, error)
}

type licenseRepository struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) LicenseRepository {
	return &licenseRepository{db: db}
}

func (l *licenseRepository) StoreKey(ctx context.Context, req domain.StoreKeyReq) (domain.License, error) {
	usr := ctx.Value("user").(domain.UserPayload)

	existing, err := gorm.G[domain.License](l.db).Where("user_id = ?", usr.UserId).First(ctx)

	if err == nil && existing.ID != 0 {
		fmt.Println("returning existing key", existing)
		return existing, nil
	}

	license := &domain.License{
		Key:      req.Key,
		Duration: req.Duration,
		UserID:   uint(usr.UserId),
		Status:   "pending",
	}

	err = gorm.G[domain.License](l.db).Create(ctx, license)
	if err != nil {
		return domain.License{}, fmt.Errorf("failed to write key", err.Error())
	}

	return *license, nil
}

func (l *licenseRepository) GetLicenseById(ctx context.Context, id string) (domain.License, error) {

	license, err := gorm.G[domain.License](l.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return domain.License{}, err
	}

	return license, nil
}

func (l *licenseRepository) GetLicenseByKey(ctx context.Context, key string) (*domain.License, error) {
	license, err := gorm.G[domain.License](l.db).Where("key = ?", key).First(ctx)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (l *licenseRepository) GetLicenseBySessionId(ctx context.Context, sessionId string) (*domain.License, error) {
	payment, err := gorm.G[domain.Payment](l.db).Where("session_id = ?", sessionId).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("payment not found for session %s: %w", sessionId, err)
	}

	license, err := gorm.G[domain.License](l.db).Where("id = ?", payment.LicenseID).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("license not found for session %s: %w", sessionId, err)
	}

	return &license, nil
}

func (l *licenseRepository) UpdateLicense(ctx context.Context, license domain.License) (domain.License, error) {

	_, err := gorm.G[domain.License](l.db).Where("key = ?", license.Key).Updates(ctx, license)

	if err != nil {
		return license, fmt.Errorf("failed to update license: %s", err.Error())
	}

	return license, nil
}
