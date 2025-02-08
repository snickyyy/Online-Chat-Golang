package domain

type User struct {
	BaseModel
	username 	string	`gorm:"unique;size:40;not null;"`
	email 		string	`gorm:"unique;size:255;not null;"`
	password 	string	`gorm:"not null"`
	description string	`gorm:"size:255;"`
	role 		string	`gorm:"size:50;not null;default:'anonymous'"`
	image 		string
}

func NewUser(username, email, password, description, role, image string) *User {
	return &User{
		username: username,
		email: email,
		password: password,
		description: description,
		role: role,
		image: image,
	}
}
