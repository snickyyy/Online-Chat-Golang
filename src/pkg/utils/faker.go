package utils

import (
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"runtime"
	"sync"
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
