package config

// status codes
const (
	ErrorCodeInvalidURL        = "INVALID_URL"
	ErrorCodeInvalidJSON       = "INVALID_JSON"
	ErrorCodeInternal          = "INTERNAL"
	ErrorCodeAlreadyExists     = "ALREADY_EXISTS"
	ErrorCodeNotFound          = "NOT_FOUND"
	ErrorBadRequest            = "BAD_REQUEST"
	ErrorCodeUnauthorized      = "UNAUTHORIZED"
	ErrorCodeForbidden         = "FORBIDDEN"
	ErrorCodeTokenExpired      = "TOKEN_EXPIRED"
	ErrorCodeNotApproved       = "NOT_APPROVED"
	ErrorCodeIncorrectPassword = "INCORRECT_PASSWORD"
	ErrorCodeInvalidToken      = "INVALID_TOKEN"
	ErrorCodeSessionLimit      = "SESSION_LIMIT"
	ErrorCodeInvalidCode       = "INVALID_CODE"
	StatusSuccess              = "SUCCESS"
	Status2faRequired          = "2FA_REQUIRED"
)

// grpc context timeout
const (
	GrpcContextTimeout = 10
)

// custom header keys
const (
	KeyLocation            = "location"
	KeySessionId           = "session_id"
	KeyUserId              = "user_id"
	KeyRoleId              = "role_id"
	KeyPlatformId          = "platform_id"
	KeyPlatformKey         = "platform_key"
	KeyTimezone            = "timezone"
	KeyDeviceName          = "device_name"
	KeyTag                 = "tag"
	KeyCompanyId           = "company_id"
	KeyCompanyUserRoleSlug = "company_user_role_slug"
	KeyCompanyUserRoleId   = "company_user_role_id"
)

// token types
const (
	AccessToken    = "access_token"
	RefreshToken   = "refresh_token"
	Verfication    = "verification"
	ForgetPassword = "forget_password"
	ResetPassword  = "reset_password"
)

// openapis
var (
	OpenApis = map[string]bool{
		"/v1/auth/login":                  true,
		"/v1/auth/logout":                 true,
		"/v1/auth/sign-up":                true,
		"/v1/auth/forget-password":        true,
		"/v1/auth/forget-password-verify": true,
		"/v1/auth/reset-password":         true,
	}
)
