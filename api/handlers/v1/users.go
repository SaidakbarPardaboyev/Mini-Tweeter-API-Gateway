package v1

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/helpers"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	usersService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/users_service"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/etc"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/logger"
	static "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/statics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// @Router /v1/users [post]
// @Summary User Create
// @Description API for creating user
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.CreateUser true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		req models.CreateUser
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
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error creating merchant") {
		return
	}

	// add user to redis
	userData := models.User{
		Id:           respCreateUser.Id,
		Name:         req.Name,
		Surname:      req.Surname,
		Bio:          req.Bio,
		ProfileImage: req.ProfileImage,
		Username:     req.Username,
		Password:     string(hashedPassword),
	}
	err = h.Redis.Set(ctx, fmt.Sprintf(static.UserCacheFormat, respCreateUser.Id), etc.MarshalJSON(userData), static.CachingExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding user to redis") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "User created sussecsfully",
		Data:    respCreateUser,
	})
}

// @Router /v1/users/{id} [get]
// @Summary Get User By ID
// @Description API for getting user by id
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "invalid id") {
			return
		}
	}

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	// get from redis
	userData, err := h.Redis.Get(ctx, fmt.Sprintf(static.UserCacheFormat, id)).Result()
	if err == redis.Nil {
		resp, err := h.services.UserService().UserService(ctx).GetUser(ctx, &usersService.GetUserRequest{
			Id: id,
		})
		if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting user by id") {
			return
		}

		err = h.Redis.Set(ctx, fmt.Sprintf(static.UserCacheFormat, id), etc.MarshalJSON(resp), static.CachingExpireTime).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding user to redis") {
			return
		}

		c.JSON(200, models.Response{
			Code:    config.StatusSuccess,
			Message: "User got sussecsfully",
			Data:    resp,
		})
		return
	} else if err != nil {
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting user from redis") {
			return
		}
	}

	var resp models.User
	err = json.Unmarshal([]byte(userData), &resp)
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshaling user data") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "User got sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/users/list [get]
// @Summary Get Users List With Pagination
// @Description API for getting users list with pagination and filters
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int true "Page number {default: 1}"
// @Param limit query int true "Limit number {default: 10}"
// @Param search query string false "Search keyword (role, login, fullname, phone, etc.)"
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sorting order (asc or desc)"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetLislUsers(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting limit from query params") {
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting page from query params") {
		return
	}

	req := &usersService.GetListUserRequest{
		Limit: int32(limit),
		Page:  int32(page),
	}

	// Extract optional search and filters
	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = &sortBy
	}
	if orderParam := c.Query("order"); orderParam != "" {
		req.Order = &orderParam
	}

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	h.log.Info("Received get list user request", logger.Any("req", req))

	resp, err := h.services.UserService().UserService(ctx).GetListUser(ctx, req)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting users list with pagination") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Users list retrieved successfully",
		Data:    resp,
	})
}

// @Router /v1/users [put]
// @Summary User Update
// @Description API for updating user
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.UpdateUser true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		req models.UpdateUser
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

	resp, err := h.services.UserService().UserService(ctx).UpdateUser(ctx, &usersService.UpdateUserRequest{
		User: &usersService.User{
			Id:           req.Id,
			Name:         req.Name,
			Surname:      req.Surname,
			Bio:          req.Bio,
			ProfileImage: req.ProfileImage,
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
		},
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error updating user") {
		return
	}

	// add user to redis
	err = h.Redis.Set(ctx, fmt.Sprintf(static.UserCacheFormat, req.Id), etc.MarshalJSON(req), static.CachingExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding user to redis") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "user updated sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/users/{id} [delete]
// @Summary Delete User
// @Description API for deleting user
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	resp, err := h.services.UserService().UserService(ctx).DeleteUser(ctx, &usersService.DeleteUserRequest{
		Id: id,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error deleting merchant") {
		return
	}

	// delete user from redis
	err = h.Redis.Del(ctx, fmt.Sprintf(static.UserCacheFormat, id)).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while deleting user from redis") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "User deleted sussecsfully",
		Data:    resp,
	})
}
