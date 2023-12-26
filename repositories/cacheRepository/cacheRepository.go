package cacheRepository

import (
	"context"
	"log"
	"time"

	"github.com/FianGumilar/email-otp-verification/repositories"
)

type cacheRepository struct {
	RepoDB repositories.Repository
}

func NewCacheRepository(RepoDB repositories.Repository) repositories.CacheRepository {
	return cacheRepository{
		RepoDB: RepoDB,
	}
}

const defineColumnCache = ` key, entry, created_at  `

// GetCache implements repositories.CacheRepoository.
func (ctx cacheRepository) GetCache(key string) ([]byte, error) {
	var val []byte
	err := ctx.RepoDB.DB.QueryRowContext(context.Background(), "SELECT entry FROM cache WHERE key = $1", key).Scan(&val)
	if err != nil {
		log.Printf("Failed to get value for key %s: %v", key, err)
		return nil, err
	}
	log.Println("Successfully retrieved value from DB")
	return val, nil

}

// SetCache implements repositories.CacheRepoository.
func (ctx cacheRepository) SetCache(key string, entry []byte) error {
	query := `INSERT INTO cache(` + defineColumnCache + `) VALUES ($1, $2, $3) ON CONFLICT (key) DO UPDATE SET entry = EXCLUDED.entry, created_at = EXCLUDED.created_at`

	_, err := ctx.RepoDB.DB.ExecContext(context.Background(), query, key, entry, time.Now())
	if err != nil {
		log.Printf("Failed to set value for key %s: %v", key, err)
		return err
	}
	log.Println("Successfully set value in DB")
	return nil
}
