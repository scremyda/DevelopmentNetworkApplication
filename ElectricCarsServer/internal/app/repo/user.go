package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) AddUser(newUser *ds.Users) error {
	result := r.db.Create(&newUser)
	if result.Error != nil {
		// Проверяем, является ли ошибка ошибкой уникального ключа
		if isDuplicateKeyError(result.Error) {
			return fmt.Errorf("login already exists")
		}
		// В противном случае, возвращаем оригинальную ошибку
		return result.Error
	}
	return nil
}

// Функция для проверки, является ли ошибка ошибкой уникального ключа
func isDuplicateKeyError(err error) bool {
	pgError, isPGError := err.(*pgconn.PgError)
	if isPGError && pgError.Code == "23505" {
		// Код "23505" является кодом ошибки уникального ключа в PostgreSQL
		return true
	}
	return false
}

func (r *Repository) SignUp(ctx context.Context, newUser ds.Users) error {
	return r.db.Create(&newUser).Error
}

func (r *Repository) GetByCredentials(ctx context.Context, user ds.Users) (ds.Users, error) {
	err := r.db.First(&user, "login = ? AND password = ?", user.Login, user.Password).Error
	return user, err
}
