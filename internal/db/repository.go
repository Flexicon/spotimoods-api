package db

import (
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/jinzhu/gorm"
)

// RepositoryProvider manages all db repositories
type RepositoryProvider struct {
	db *gorm.DB
}

// NewRepositoryProvider constructor
func NewRepositoryProvider(db *gorm.DB) internal.RepositoryProvider {
	return &RepositoryProvider{db: db}
}

// User returns a new UserRepository
func (p *RepositoryProvider) User() internal.UserRepository {
	return &UserRepository{db: p.db}
}
