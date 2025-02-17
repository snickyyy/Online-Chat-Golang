package repositories

import domain "libs/src/internal/domain/models"

type UserRepository struct {
	BasePostgresRepository[domain.User]
}
