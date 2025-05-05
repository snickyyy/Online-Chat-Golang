package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type FakeUser struct {
	Username    string
	Email       string
	Password    string
	Description string
	IsActive    bool
	Role        byte
	Image       string
}

type FakeChat struct {
	Title       string
	Description string
	OwnerID     int64
}

type FakeMessage struct {
	ID primitive.ObjectID

	ChatID   int64
	SenderID int64
	Content  string

	IsRead    bool
	IsUpdated bool
	IsDeleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFakeUser() *FakeUser {
	fake := gofakeit.New(uint64(rand.Int63()))
	return &FakeUser{
		Username:    fake.Username(),
		Email:       fake.Email(),
		Password:    fake.Password(true, true, true, false, false, 10),
		Description: fake.Sentence(10),
		IsActive:    fake.Bool(),
		Role:        byte(fake.Number(0, 2)),
		Image:       fake.Sentence(50),
	}
}

func NewFakeChat(ownerId int64) *FakeChat {
	return &FakeChat{
		Title:       gofakeit.Username(),
		Description: gofakeit.Sentence(8),
		OwnerID:     ownerId,
	}
}

func NewFakeMessage(chatId, senderId int64) *FakeMessage {
	createdAt := time.Unix(rand.Int63n(time.Now().Unix()-94608000)+94608000, 0)
	updatedAt := createdAt.Add(time.Duration(rand.Int63n(int64(time.Since(createdAt)))))

	return &FakeMessage{
		ID:        primitive.NewObjectID(),
		ChatID:    chatId,
		SenderID:  senderId,
		Content:   gofakeit.Sentence(10),
		IsRead:    gofakeit.Bool(),
		IsUpdated: gofakeit.Bool(),
		IsDeleted: gofakeit.Bool(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func GenerateInParallel[T any](totalCount int, action func(count int) []T) []T {
	result := []T{}
	countPerWorker := int(totalCount / runtime.NumCPU())

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := action(countPerWorker)
			mutex.Lock()
			result = append(result, data...)
			mutex.Unlock()
		}()
	}
	wg.Wait()
	return result
}
