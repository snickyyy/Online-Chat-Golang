package repositories

import domain "libs/src/internal/domain/models"

// import (
// 	"context"
// 	domain "libs/src/internal/domain/models"
// 	"libs/src/settings"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

type ChatRepository struct {
	BaseRepository[domain.Chat]
}
