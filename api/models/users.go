package models

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateUser struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type UpdateUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}
