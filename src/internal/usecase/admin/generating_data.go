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
	"sync"
	"time"
)

type DataGenerator struct {
	App                  *settings.App
	UserRepository       repositories.IUserRepository
	ChatRepository       repositories.IChatRepository
	ChatMemberRepository repositories.IChatMemberRepository
	MessageRepository    repositories.IMessageRepository
}

func NewDataGenerator(app *settings.App) *DataGenerator {
	return &DataGenerator{
		App:                  app,
		UserRepository:       repositories.NewUserRepository(app),
		ChatRepository:       repositories.NewChatRepository(app),
		ChatMemberRepository: repositories.NewChatMemberRepository(app),
		MessageRepository:    repositories.NewMessageRepository(app),
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

func (dg *DataGenerator) GenerateChatMembers(caller dto.UserDTO, count int) error {
	if caller.Role < enums.ADMIN || !caller.IsActive {
		return usecase_errors.PermissionError{Msg: "You are not allowed to perform this action"}
	}

	chats, err := dg.ChatRepository.GetAll(dg.App.Ctx)
	if err != nil {
		return err
	}
	if len(chats) < 1 {
		return usecase_errors.NotFoundError{Msg: "No chats found"}
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
		members := make([]domain.ChatMember, count)

		for i := 0; i < count; i++ {
			chatId := chats[rand.IntN(len(chats))].ID
			userId := users[rand.IntN(len(users))].ID
			memberRole := rand.IntN(3)
			members[i] = domain.ChatMember{
				ChatID:     chatId,
				UserID:     userId,
				MemberRole: byte(memberRole),
			}
		}

		fmt.Printf("Generated %d chat members in %dms", count, int(time.Since(start).Milliseconds()))

		start = time.Now()
		err := dg.ChatMemberRepository.ManyToCreate(dg.App.Ctx, members)
		fmt.Printf("Saved %d chat members in %dms", count, int(time.Since(start).Milliseconds()))
		if err != nil {
			dg.App.Logger.Error(fmt.Sprintf("Error generating chat members: %v", err))
		}
	}()

	return nil
}

func (dg *DataGenerator) GenerateMessages(caller dto.UserDTO, count int) error {
	if caller.Role < enums.ADMIN || !caller.IsActive {
		return usecase_errors.PermissionError{Msg: "You are not allowed to perform this action"}
	}

	wg := &sync.WaitGroup{}

	var chats []domain.Chat
	var errFromChat error
	wg.Add(2)
	go func() {
		defer wg.Done()
		chats, errFromChat = dg.ChatRepository.GetAll(dg.App.Ctx)
	}()

	var users []domain.User
	var errFromUser error

	go func() {
		defer wg.Done()
		users, errFromUser = dg.UserRepository.GetAll(dg.App.Ctx)
	}()

	wg.Wait()

	if errFromUser != nil {
		return errFromUser
	}
	if len(users) < 1 {
		return usecase_errors.NotFoundError{Msg: "No users found"}
	}

	if errFromChat != nil {
		return errFromChat
	}
	if len(chats) < 1 {
		return usecase_errors.NotFoundError{Msg: "No chats found"}
	}

	go func() {
		action := func(count int) []domain.Message {
			messages := make([]domain.Message, count)
			for i := 0; i < count; i++ {
				chatId := chats[rand.IntN(len(chats))].ID
				senderId := users[rand.IntN(len(users))].ID
				message := utils.NewFakeMessage(chatId, senderId)
				messages[i] = domain.Message{
					BaseMongo: domain.BaseMongo{
						Id:        message.ID,
						CreatedAt: message.CreatedAt,
						UpdatedAt: message.UpdatedAt,
					},
					SenderId:  senderId,
					ChatId:    chatId,
					Content:   message.Content,
					IsRead:    message.IsRead,
					IsUpdated: message.IsUpdated,
					IsDeleted: message.IsDeleted,
				}
			}
			return messages
		}
		start := time.Now()
		messages := utils.GenerateInParallel[domain.Message](count, action)
		fmt.Printf("Generated %d messages in %dms", count, int(time.Since(start).Milliseconds()))
		start = time.Now()
		err := dg.MessageRepository.CreateMany(dg.App.Ctx, messages)
		fmt.Printf("Saved %d messages in %dms", count, int(time.Since(start).Milliseconds()))
		if err != nil {
			dg.App.Logger.Error(fmt.Sprintf("Error generating messages: %v", err))
		}
	}()

	return nil
}
