package models

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	UserId string `json:"user_id"`
}

type ResetPasswordRequest struct {
	TokenKey string `json:"token"`
	Password string `json:"password"`
}

type ResetTokenRequest struct {
	TokenKey string `json:"token"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type SignUpResponse struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshTokenRequest struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type TokenGenerationResponse struct {
	CompanyUserToken string `json:"company_user_token"`
}

type VerifyEmail struct {
	Verify bool `json:"verify"`
}

type GetMenu struct {
	Token string `json:"token"`
}

type Menu struct {
	Id        string `json:"id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	ParentId  string `json:"parent_id"`
	CanCreate bool   `json:"can_create"`
	CanUpdate bool   `json:"can_update"`
	CanGet    bool   `json:"can_get"`
	CanDelete bool   `json:"can_delete"`
	Icon      string `json:"icon"`
}

type UploadResponse struct {
	OriginalURL string `json:"filename,omitempty"`
	Key         string `json:"key"`
}

type SignUpLink struct {
	Link string `json:"link"`
}
