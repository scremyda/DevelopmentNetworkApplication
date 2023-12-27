package repo

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetUserInfo(r *Repository, id uint) (ds.Users, error) {
	var user ds.Users

	result := r.db.Where("id = ?", id).Find(&user)

	if result.Error != nil {
		return ds.Users{}, result.Error
	}

	return user, nil
}

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
