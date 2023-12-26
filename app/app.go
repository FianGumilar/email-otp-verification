package app

import (
	"database/sql"

	"github.com/FianGumilar/email-otp-verification/repositories"
	"github.com/FianGumilar/email-otp-verification/repositories/cacheRepository"
	"github.com/FianGumilar/email-otp-verification/repositories/userRepository"
	"github.com/FianGumilar/email-otp-verification/services"
)

func SetupApp(DB *sql.DB, repo repositories.Repository) services.UsecaseService {

	// Repository
	userRepo := userRepository.NewUserRepository(repo)
	cacheRepo := cacheRepository.NewCacheRepository(repo)

	// Services
	usecaseSvc := services.NewUsecaseService(DB, userRepo, cacheRepo)

	return usecaseSvc
}
