package repositories

import (
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
)

//go:generate mockery --name=IUserRepository --dir=. --output=../mocks --with-expecter
type IUserRepository interface {
	GetByUsername(username string) (domain.User, error)
	IBasePostgresRepository[domain.User]
}

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

func (r *UserRepository) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	err := r.Db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return domain.User{}, parsePgError(err)
	}
	return user, nil
}
