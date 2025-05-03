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
	"math/rand/v2"
	"strconv"
	"time"
)

type DataGenerator struct {
	App            *settings.App
	UserRepository repositories.IUserRepository
	ChatRepository repositories.IChatRepository
}

func NewDataGenerator(app *settings.App) *DataGenerator {
	return &DataGenerator{
		App:            app,
		UserRepository: repositories.NewUserRepository(app),
		ChatRepository: repositories.NewChatRepository(app),
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

func (dg *DataGenerator) GenerateChats(caller dto.UserDTO, count int) error {
	if caller.Role < enums.ADMIN || !caller.IsActive {
		return usecase_errors.PermissionError{Msg: "You are not allowed to perform this action"}
	}

	users, err := dg.UserRepository.GetAll(dg.App.Ctx)
	if err != nil {
		return err
	}
	if len(users) < 1 {
		return usecase_errors.NotFoundError{Msg: "No users found"}
	}

	go func() {
		start := time.Now()
		funcForGenerateAndParseChat := func(count int) []domain.Chat {
			chats := make([]domain.Chat, count)
			for i := 0; i < count; i++ {
				ownerId := users[rand.IntN(len(users))].ID
				chat := utils.NewFakeChat(ownerId)
				chats[i] = domain.Chat{
					Title:       chat.Title + strconv.Itoa(i),
					Description: chat.Description,
					OwnerID:     ownerId,
				}
			}
			return chats
		}
		chats := utils.GenerateInParallel[domain.Chat](count, funcForGenerateAndParseChat)
		fmt.Printf("Generated %d chats in %dms", count, int(time.Since(start).Milliseconds()))

		start = time.Now()
		err := dg.ChatRepository.ManyToCreate(dg.App.Ctx, chats)
		fmt.Printf("Saved %d chats in %dms", count, int(time.Since(start).Milliseconds()))
		if err != nil {
			dg.App.Logger.Error(fmt.Sprintf("Error generating chats: %v", err))
		}
	}()

	return nil
}
