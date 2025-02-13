package domain


type Chat struct {
	BaseMongo
	OwnerId		int64 	`bson:"owner_id" json:"owner_id,omitempty"`
	Title		string 	`bson:"title" json:"title"`
	Description	string 	`bson:"description" json:"description"`
	Members		[]int64 `bson:"members"`
}

type Message struct {
	BaseMongo
	SenderId 	string	`bson:"sender_id" json:"sender_id"`
	ChatId		int64   `bson:"chat_id" json:"chat_id"`
	MessageType	string 	`bson:"message_type" json:"message_type"`
	Content     string  `bson:"content" json:"content"`
	IsUpdated	bool 	`bson:"is_updated" json:"is_updated"`
	IsDeleted	bool `bson:"is_deleted" json:"is_deleted"`
}
