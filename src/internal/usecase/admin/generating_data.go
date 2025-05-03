package services

import (
	"fmt"
	"libs/src/internal/domain/enums"
	domain "libs/src/internal/domain/models"
	"libs/src/internal/dto"
	"libs/src/internal/repositories"
	usecase_errors "libs/src/internal/usecase/errors"
	"libs/src/pkg/utils"
	"libs/src/settings"
	"strconv"
	"time"
)

type DataGenerator struct {
	App            *settings.App
	UserRepository repositories.IUserRepository
}

func NewDataGenerator(app *settings.App) *DataGenerator {
	return &DataGenerator{
		App:            app,
		UserRepository: repositories.NewUserRepository(app),
	}
}

func (dg *DataGenerator) GenerateUsers(caller dto.UserDTO, count int) error {
	if caller.Role < enums.ADMIN || !caller.IsActive {
		return usecase_errors.PermissionError{Msg: "You are not allowed to perform this action"}
	}

	go func() {
		start := time.Now()
		funcForGenerateAndParseUser := func(count int) []domain.User {
			users := make([]domain.User, count)
			for i := 0; i < count; i++ {
				user := utils.NewFakeUser()
				users[i] = domain.User{
					Username:    user.Username + strconv.Itoa(i),
					Email:       user.Email + strconv.Itoa(i),
					Password:    user.Password,
					Description: user.Description,
					IsActive:    user.IsActive,
					Role:        user.Role,
					Image:       user.Image,
				}
			}
			return users
		}
		users := utils.GenerateInParallel[domain.User](count, funcForGenerateAndParseUser)
		fmt.Printf("Generated %d users in %dms", count, int(time.Since(start).Milliseconds()))

		start = time.Now()
		err := dg.UserRepository.ManyToCreate(dg.App.Ctx, users)
		fmt.Printf("Saved %d users in %dms", count, int(time.Since(start).Milliseconds()))
		if err != nil {
			dg.App.Logger.Error(fmt.Sprintf("Error generating users: %v", err))
		}
	}()

	return nil
}
