package repositories

import (
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

type UserRepository struct {
	BasePostgresRepository[domain.User]
}

func NewUserRepository(app *settings.App) *UserRepository {
	return &UserRepository{
		BasePostgresRepository: BasePostgresRepository[domain.User]{
			Model: domain.User{},
			Db:    app.DB,
		},
	}
}
