package repository

import (
	"context"
	"fmt"
	"license-server/internals/domain"

	"gorm.io/gorm"
)

type LicenseRepository interface {
	StoreKey(ctx context.Context, key string) (domain.License, error)
	GetLicenseById(ctx context.Context, id string) (domain.License, error)
	GetLicenseByKey(ctx context.Context, key string) (*domain.License, error)
	UpdateLicense(ctx context.Context, license domain.License) (domain.License, error)
}

type licenseRepository struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) LicenseRepository {
	return &licenseRepository{db: db}
}

func (l *licenseRepository) StoreKey(ctx context.Context, key string) (domain.License, error) {
	usr := ctx.Value("user").(domain.UserPayload)

	existing, err := gorm.G[domain.License](l.db).Where("user_id = ?", usr.UserId).First(ctx)

	if err == nil && existing.ID != 0 {
		return existing, nil
	}

	license := &domain.License{
		Key:    key,
		UserID: uint(usr.UserId),
		Status: "pending",
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

func (l *licenseRepository) UpdateLicense(ctx context.Context, license domain.License) (domain.License, error) {

	_, err := gorm.G[domain.License](l.db).Updates(ctx, domain.License{
		MachineID: license.MachineID,
		Status:    license.MachineID,
	})

	if err != nil {
		return license, fmt.Errorf("failed to update license: %s", err.Error())
	}

	return license, nil
}
