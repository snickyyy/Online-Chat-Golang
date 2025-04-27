package repositories

import (
	"context"
	domain "libs/src/internal/domain/models"
	"libs/src/settings"
	"time"
)

//go:generate mockery --name=IUserRepository --dir=. --output=../mocks --with-expecter
type IUserRepository interface {
	GetByUsername(Ctx context.Context, username string) (domain.User, error)
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

func (r *UserRepository) GetByUsername(Ctx context.Context, username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(Ctx, time.Duration(settings.AppVar.Config.Timeout.Postgres.Medium)*time.Millisecond)
	defer cancel()

	var user domain.User
	err := r.Db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return domain.User{}, parsePgError(err)
	}
	return user, nil
}
