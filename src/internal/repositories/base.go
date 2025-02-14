package repositories

import "go.mongodb.org/mongo-driver/mongo"

type BaseRepository struct {
	Db *mongo.Database
	
}
