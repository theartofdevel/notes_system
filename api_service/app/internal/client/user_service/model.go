package user_service

type User struct {
	UUID     string `json:"uuid" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password,omitempty"`
}

type SigninUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password,omitempty"`
}

type CreateUserDTO struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type UpdateUserDTO struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}
