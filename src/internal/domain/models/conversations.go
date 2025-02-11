package domain


type Chat struct {
	BaseMongo
	OwnerId		int64 	`bson:"owner_id,omitempty" json:"owner_id,omitempty"`
	Title		string 	`bson:"title,omitempty" json:"title"`
	Description	string 	`bson:"description" json:"description"`
	Members		[]int64 `bson:"members"`
}
