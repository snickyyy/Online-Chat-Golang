package domain

import (
	"reflect"
	"time"
)

type BaseModel struct {
	ID        int64 		`gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseMongo struct {
	Id 	  		string 		`bson:"_id,omitempty" json:"id"`
	CreatedAt 	time.Time	`bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt 	time.Time	`bson:"updated_at" json:"updated_at,omitempty"`
	DeleteAt 	time.Time	`bson:"delete_at" json:"delete_at,omitempty"`
}

func (bm BaseMongo) Mapping() map[string]interface{} {
	result := make(map[string]interface{})
	type_ := reflect.TypeOf(bm)
	value := reflect.ValueOf(bm)

	for i := 0; i < value.NumField(); i++ {
		result[type_.Field(i).Name] = value.Field(i)
	}
	
	return result
}
