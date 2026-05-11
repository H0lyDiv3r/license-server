package store

import (
	"client/domain"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".config", "secure_desktop")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, "app.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.License{})

	return &Store{
		Db: db,
	}, nil
}

func (s *Store) StoreLicence(ctx context.Context, license domain.License) (*domain.License, error) {
	fmt.Println("going to store the license here", license)

	existing, err := gorm.G[domain.License](s.Db).Where("key = ?", license.Key).First(ctx)
	if err == nil && existing.ID != 0 {
		return &existing, nil
	}

	err = gorm.G[domain.License](s.Db).Create(ctx, &license)
	if err != nil {
		return nil, err
	}

	return &license, nil
}

func (s *Store) UpdateLicense(ctx context.Context, license domain.License) (*domain.License, error) {

	existing, err := gorm.G[domain.License](s.Db).Where("key = ?", license.Key).First(ctx)
	if err != nil {
		return nil, err
	}

	if existing.ID == 0 {
		return nil, fmt.Errorf("license doesnt exist")
	}

	_, err = gorm.G[domain.License](s.Db).Where("key = ?", license.Key).Updates(ctx, license)
	if err != nil {
		return nil, err
	}

	return &license, nil
}
