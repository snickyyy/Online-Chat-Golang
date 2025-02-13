package domain

type User struct {
	BaseModel
	Username 	string	`gorm:"unique;size:40;not null;"`
	Email 		string	`gorm:"unique;size:255;not null;"`
	Password 	string	`gorm:"not null"`
	Description string	`gorm:"size:255;"`
	Role 		string	`gorm:"size:50;not null;default:'anonymous'"`
	Image 		string
}

func NewUser(username, email, password, description, role, image string) *User {
	return &User{
		Username: username,
		Email: email,
		Password: password,
		Description: description,
		Role: role,
		Image: image,
	}
}
