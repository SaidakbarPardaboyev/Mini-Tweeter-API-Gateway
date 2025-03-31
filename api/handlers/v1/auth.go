package v1

import (
	"fmt"
	"time"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/helpers"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	userService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/users_service"
	usersService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/users_service"
	etc "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/etc"
	jwt "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/jwt"
	static "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/statics"
	"github.com/gin-gonic/gin"

	// "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// @Router /v1/auth/sign-up [post]
// @Summary User Sign Up
// @Description API for user sign up
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body models.CreateUser true "request"
// @Success 200 {object} models.Response
// @Failure default {object} models.Response
func (h *handlerV1) SignUp(c *gin.Context) {
	var (
		req              models.CreateUser
		userData         models.User
		currentTimestamp = time.Now()
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	// hashing password
	hashedPassword, err := etc.GeneratePasswordHash(req.Password)
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while hashing password") {
		return
	}

	respCreateUser, err := h.services.UserService().UserService(ctx).CreateUser(ctx, &usersService.CreateUserRequest{
		User: &usersService.User{
			Id:           "",
			Name:         req.Name,
			Surname:      req.Surname,
			Bio:          req.Bio,
			ProfileImage: req.ProfileImage,
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
		},
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error creating user") {
		return
	}
	userData = models.User{
		Id:           respCreateUser.Id,
		Name:         req.Name,
		Surname:      req.Surname,
		Bio:          req.Bio,
		ProfileImage: req.ProfileImage,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    currentTimestamp.Format(time.RFC3339),
	}

	// add user to redis
	err = h.RedisClient().Set(ctx, fmt.Sprintf(static.UserCacheFormat, userData.Id), etc.MarshalJSON(userData), static.CachingExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while storing user to redis") {
		return
	}

	// generating access token
	accessToken, err := jwt.GenerateAccessToken(&userData)
	if helpers.HandleInternalWithMessage(c, h.log, err, "error generating access token") {
		return
	}

	// generating refresh token
	refreshToken, err := jwt.GenerateRefreshToken(&userData)
	if helpers.HandleInternalWithMessage(c, h.log, err, "error generating refresh token") {
		return
	}

	// Validate Access Token
	_, _, err = jwt.ValidateToken(accessToken, static.AccessTokenSecret)
	if helpers.HandleInternalWithMessage(c, h.log, err, "error validating access token") {
		return
	}

	// store refresh token to redis
	err = h.RedisClient().Set(ctx, fmt.Sprintf(static.RefreshTokenKeyFormat, userData.Id), refreshToken, static.RefreshTokenExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while storing refresh token to redis") {
		return
	}

	fmt.Println("refresh token: ", h.RedisClient().Get(ctx, fmt.Sprintf("%s user refresh token", userData.Id)).Val())

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "signed up successfully",
		Data: models.SignUpResponse{
			UserId:       userData.Id,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(static.AccessTokenExpireTime.Seconds()),
		},
	})
}

// @Router /v1/auth/login [post]
// @Summary User Login
// @Description API for user login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body models.LoginRequest true "request"
// @Success 200 {object} models.Response
// @Failure default {object} models.Response
func (h *handlerV1) Login(c *gin.Context) {
	var (
		req models.LoginRequest
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	// getting user info from db
	userResp, err := h.services.UserService().UserService(ctx).GetUserWithLogin(ctx, &userService.GetUserWithUsernameRequest{
		Username: req.Username,
	})
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting user info") {
		return
	}

	// checking user is found
	if userResp.Id == "" {
		helpers.HandleBadRequestErrWithMessage(c, h.log, err, "User not found")
		return
	}

	// checking password is correct
	err = bcrypt.CompareHashAndPassword([]byte(userResp.PasswordHash), []byte(req.Password))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while comparing password") {
		return
	}

	userData := models.User{
		Id:           userResp.Id,
		Name:         userResp.Name,
		Surname:      userResp.Surname,
		Bio:          userResp.Bio,
		ProfileImage: userResp.ProfileImage,
		Username:     userResp.Username,
		PasswordHash: userResp.PasswordHash,
		CreatedAt:    userResp.CreatedAt,
		UpdatedAt:    userResp.UpdatedAt,
	}

	// generating access token
	accessToken, err := jwt.GenerateAccessToken(&userData)
	if helpers.HandleInternalWithMessage(c, h.log, err, "error generating access token") {
		return
	}

	// generating refresh token
	refreshToken, err := jwt.GenerateRefreshToken(&userData)
	if helpers.HandleInternalWithMessage(c, h.log, err, "error generating refresh token") {
		return
	}

	// Validate Access Token
	_, _, err = jwt.ValidateToken(accessToken, static.AccessTokenSecret)
	if helpers.HandleInternalWithMessage(c, h.log, err, "error validating access token") {
		return
	}

	// store refresh token to redis
	err = h.RedisClient().Set(ctx, fmt.Sprintf(static.RefreshTokenKeyFormat, userData.Id), refreshToken, static.RefreshTokenExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while storing refresh token to redis") {
		return
	}

	fmt.Println("refresh token: ", h.RedisClient().Get(ctx, fmt.Sprintf("%s user refresh token", userData.Id)).Val())

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Logged in successfully",
		Data: models.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(static.AccessTokenExpireTime.Seconds()),
		},
	})
}

// @Router /v1/auth/logout [post]
// @Summary Logout Endpoint
// @Description API for user to logout - actually used to delete session, `Authorization` required
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body models.LogoutRequest true "request"
// @Success 200 {object} models.Response
// @Failure default {object} models.Response
func (h *handlerV1) Logout(c *gin.Context) {
	var (
		req models.LogoutRequest
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	err = h.Redis.Del(ctx, fmt.Sprintf(static.RefreshTokenKeyFormat, req.UserId)).Err()
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error while deleting refresh token") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Logout success",
		Data:    nil,
	})
}

// // @Router /v1/auth/refresh-token [get]
// // @Summary Refresh Token
// // @Description API for reset token
// // @Tags auth
// // @Accept  json
// // @Produce  json
// // @Param request body models.ResetTokenResponse true "request"
// // @Success 200 {object} models.Response
// // @Failure default {object} models.Response
// func (h *handlerV1) RefreshToken(c *gin.Context) {
// 	var (
// 		req = models.RefreshTokenRequest{}
// 	)

// 	err := c.ShouldBindJSON(&req)
// 	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
// 		return
// 	}
// 	ctx, cancel := etc.NewTimoutContext(c)
// 	defer cancel()

// 	refreshToken, err := h.Redis.Get(ctx, fmt.Sprintf(static.RefreshTokenKeyFormat, req.UserId)).Result()
// 	if err != redis.Nil {
// 		if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Error: refresh token not found") {
// 			return
// 		}
// 	} else if err != nil {
// 		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting refresh token") {
// 			return
// 		}
// 	}

// 	// checking refresh token is correct
// 	if refreshToken != req.RefreshToken {
// 		if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid refresh token") {
// 			return
// 		}
// 	}

// 	// Simulating token refresh
// 	newAccessToken, err := jwt.RefreshToken(req.RefreshToken)
// 	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while refreshing token") {
// 		return
// 	}

// 	c.JSON(200, models.Response{
// 		Code:    config.StatusSuccess,
// 		Message: "refreshed token successfully",
// 		Data: models.LoginResponse{
// 			AccessToken:  newAccessToken,
// 			RefreshToken: refreshToken,
// 			ExpiresIn:    int64(static.AccessTokenExpireTime.Seconds()),
// 		},
// 	})
// }
