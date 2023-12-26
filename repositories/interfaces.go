package repositories

import "github.com/FianGumilar/email-otp-verification/models"

type UserRepository interface {
	AddUser(user models.User) (id int64, err error)
	EditUser(user models.User) (err error)
	RemoveUser(userId int64) (err error)
	IsUserExistsById(userId int64) (user models.User, exists bool)
	GetUserListByIndex(user models.User) (result []models.User, err error)
	IsEmailExistsByIndex(email string) (user models.User, exists bool)
	IsUsernameExistsByIndex(username string) (user models.User, exists bool)
}

type CacheRepository interface {
	GetCache(key string) ([]byte, error)
	SetCache(key string, entry []byte) error
}
