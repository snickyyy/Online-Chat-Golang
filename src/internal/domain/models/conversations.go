package domain


type Chat struct {
	BaseMongo
	OwnerId		int64 	`bson:"owner_id" json:"owner_id,omitempty"`
	Title		string 	`bson:"title" json:"title"`
	Description	string 	`bson:"description" json:"description"`
	Members		[]int64 `bson:"members"`
}
