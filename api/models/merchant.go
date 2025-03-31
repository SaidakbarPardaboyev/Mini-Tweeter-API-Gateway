package models

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
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

type UpdateMerchant struct {
	Id                string `json:"id"`
	Login             string `json:"login"`
	Password          string `json:"password"`
	Photo             string `json:"photo"`
	Fullname          string `json:"fullname"`
	Phone             string `json:"phone"`
	PassportId        string `json:"passport_id"`
	WorkplaceName     string `json:"workplace_name"`
	WorkspaceAddress  string `json:"workspace_address"`
	WorkplaceLocation string `json:"workplace_location"`
	Status            string `json:"status"`
}
