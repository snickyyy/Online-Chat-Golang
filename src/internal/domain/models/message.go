package domain

type Message struct {
	BaseMongo
	SenderId  int64  `bson:"sender_id" json:"sender_id"`
	ChatId    string `bson:"chat_id" json:"chat_id"`
	Content   string `bson:"content" json:"content"`
	IsRead    bool   `bson:"is_read" json:"is_read"`
	IsUpdated bool   `bson:"is_updated" json:"is_updated"`
	IsDeleted bool   `bson:"is_deleted" json:"is_deleted"`
}
