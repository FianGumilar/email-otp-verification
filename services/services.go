package services

import (
	"database/sql"

	"github.com/FianGumilar/email-otp-verification/repositories"
)

type UsecaseService struct {
	RepoDB    *sql.DB
	UserRepo  repositories.UserRepository
	CacheRepo repositories.CacheRepository
}

func NewUsecaseService(repoDB *sql.DB,
	userRepo repositories.UserRepository,
	cacheRepo repositories.CacheRepository,
) UsecaseService {
	return UsecaseService{
		RepoDB:    repoDB,
		UserRepo:  userRepo,
		CacheRepo: cacheRepo,
	}
}
