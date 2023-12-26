package userRepository

import (
	"database/sql"

	"github.com/FianGumilar/email-otp-verification/constans"
	. "github.com/FianGumilar/email-otp-verification/utils"

	"github.com/FianGumilar/email-otp-verification/models"
	"github.com/FianGumilar/email-otp-verification/repositories"
)

type userRepository struct {
	RepoDB repositories.Repository
}

func NewUserRepository(repoDB repositories.Repository) repositories.UserRepository {
	return userRepository{
		RepoDB: repoDB,
	}
}

const defineColumnUsers = ` full_name, phone, email, username, password, email_verified_at, created_at, updated_at `

// AddUser implements repositories.UserRepository.
func (ctx userRepository) AddUser(user models.User) (id int64, err error) {
	var args []interface{}

	query := `INSERT INTO users (` + defineColumnUsers + `) 
			VALUES (` + QueryFill(defineColumnUsers) + `) RETURNING id`

	args = append(args, user.FullName, user.Phone, user.Email, user.Username, user.Password, user.EmailVerifiedAtDB, user.EmailVerifiedAt, user.CreatedAt, user.UpdatedAt)

	query = ReplaceSQL(query, "?")
	err = ctx.RepoDB.DB.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return
	}
	return
}

// EditUser implements repositories.UserRepository.
func (ctx userRepository) EditUser(user models.User) (err error) {
	var args []interface{}

	user.EmailVerifiedAtDB = sql.NullTime{
		Time:  user.EmailVerifiedAt,
		Valid: true,
	}

	query := ` UPDATE users
	SET full_name = ?, phone = ?, email = ?, username = ?, password = ?, email_verified_at = ?, updated_at = ?,
	WHERE id = ?
	`
	args = append(args, user.FullName, user.Phone, user.Email, user.Username, user.Password, user.EmailVerifiedAt, user.UpdatedAt, user.ID)

	_, err = ctx.RepoDB.DB.Exec(query, args...)
	if err != nil {
		return
	}
	return
}

// GetUserListByIndex implements repositories.UserRepository.
func (ctx userRepository) GetUserListByIndex(user models.User) (result []models.User, err error) {
	var args []interface{}

	query := `SELECT id` + defineColumnUsers + `FROM users`

	if user.FullName != constans.EMPTY_VALUE {
		query += ` AND agent_name ILIKE ? || '%'  `
		args = append(args, user.FullName)
	}

	if user.Phone != constans.EMPTY_VALUE {
		query += ` AND phone ILIKE ? || '%'  `
		args = append(args, user.Phone)
	}

	if user.Email != constans.EMPTY_VALUE {
		query += ` AND email ILIKE ? || '%'  `
		args = append(args, user.Email)
	}

	if user.Username != constans.EMPTY_VALUE {
		query += ` AND username ILIKE ? || '%'  `
		args = append(args, user.Username)
	}

	if user.Password != constans.EMPTY_VALUE {
		query += ` AND password ILIKE ? || '%'  `
		args = append(args, user.Password)
	}

	query = ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.Query(query, args...)
	if err != nil {
		return
	}

	defer rows.Close()
	data, err := userDto(rows)
	if len(data) == constans.EMPTY_VALUE_INT {
		return result, nil
	}
	return data, err
}

// IsUserExistsById implements repositories.UserRepository.
func (ctx userRepository) IsUserExistsById(userId int64) (user models.User, exists bool) {
	query := `SELECT id, ` + defineColumnUsers + `FROM users WHERE id = $1`
	rows, err := ctx.RepoDB.DB.Query(query, userId)
	if err != nil {
		return user, constans.FALSE_VALUE
	}

	defer rows.Close()
	data, _ := userDto(rows)
	if len(data) == constans.EMPTY_VALUE_INT {
		return user, constans.FALSE_VALUE
	}
	return data[0], constans.TRUE_VALUE
}

// IsEmailExistsByIndex implements repositories.UserRepository.
func (ctx userRepository) IsEmailExistsByIndex(email string) (user models.User, exists bool) {
	query := `SELECT id, ` + defineColumnUsers + ` FROM users WHERE email = $1`

	rows, err := ctx.RepoDB.DB.Query(query, email)
	if err != nil {
		return user, constans.FALSE_VALUE
	}

	defer rows.Close()
	data, _ := userDto(rows)
	if len(data) == constans.EMPTY_VALUE_INT {
		return user, constans.FALSE_VALUE
	}
	return data[0], constans.TRUE_VALUE
}

// IsUsernameExistsByIndex implements repositories.UserRepository.
func (ctx userRepository) IsUsernameExistsByIndex(username string) (user models.User, exists bool) {
	query := `SELECT id, ` + defineColumnUsers + ` FROM users WHERE username = $1`

	rows, err := ctx.RepoDB.DB.Query(query, username)
	if err != nil {
		return user, constans.FALSE_VALUE
	}

	defer rows.Close()
	data, _ := userDto(rows)
	if len(data) == constans.EMPTY_VALUE_INT {
		return user, constans.FALSE_VALUE
	}
	return data[0], constans.TRUE_VALUE
}

// RemoveUser implements repositories.UserRepository.
func (ctx userRepository) RemoveUser(userId int64) (err error) {
	query := `DELETE FROM users WHERE id = ?`

	query = ReplaceSQL(query, "?")
	_, err = ctx.RepoDB.DB.Exec(query, userId)
	if err != nil {
		return err
	}
	return
}

// User DTO
func userDto(rows *sql.Rows) (result []models.User, err error) {
	for rows.Next() {
		var val models.User
		err = rows.Scan(&val.ID, &val.FullName, &val.Phone, &val.Email, &val.Username, &val.Password,
			&val.EmailVerifiedAt, &val.EmailVerifiedAtDB, &val.CreatedAt, &val.UpdatedAt)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}
