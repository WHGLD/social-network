package models

const (
	GenderMale   = "male"
	GenderFemale = "female"
)

type User struct {
	ID           string  `json:"user_id"`
	Password     string  `json:"password"` // нужен для парсинга в регистр методе
	PasswordHash string  `json:"password_hash"`
	FirstName    string  `json:"first_name"`
	SecondName   *string `json:"second_name"`
	Birthday     *string `json:"birthday"`
	Sex          *string `json:"sex"`
	Biography    *string `json:"biography"`
	City         *string `json:"city"`
}

type UserResponse struct {
	ID         string  `json:"user_id"`
	FirstName  string  `json:"first_name"`
	SecondName *string `json:"second_name"`
	Birthday   *string `json:"birthday"`
	Sex        *string `json:"sex"`
	Biography  *string `json:"biography"`
	City       *string `json:"city"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		Birthday:   u.Birthday,
		Sex:        u.Sex,
		Biography:  u.Biography,
		City:       u.City,
	}
}

func UsersTransform(users []User) (usersResponse []UserResponse) {
	for _, user := range users {
		usersResponse = append(usersResponse, user.ToResponse())
	}
	return
}
