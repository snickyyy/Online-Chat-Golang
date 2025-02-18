package filters

type ChatFilterNonStrict struct {
	OwnerId		int64 	`bson:"owner_id,omitempty" json:"owner_id,omitempty"`
	Title		string 	`bson:"title,omitempty" json:"title"`
	Description	string 	`bson:"description,omitempty" json:"description"`
	Members		[]int64 `bson:"members,omitempty"`
}

type ChatFilterStrict struct {
	OwnerId		int64 	`bson:"owner_id" json:"owner_id,omitempty"`
	Title		string 	`bson:"title" json:"title"`
	Description	string 	`bson:"description" json:"description"`
	Members		[]int64 `bson:"members"`
}

func (ChatFilterStrict) ApplyFilter(){}

func (ChatFilterNonStrict) ApplyFilter(){}
