package postgres

import (
	"amani-engine/internal/domain"
)

// UserRepo — структура для работы с пользователями в БД
type UserRepo struct {
	// Здесь позже будет подключение к БД
}

// GetUserByID — пример функции, которую мы напишем позже
func (r *UserRepo) GetUserByID(id int) (*domain.User, error) {
	return nil, nil
}
