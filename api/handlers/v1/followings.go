package v1

import (
	"fmt"
	"strconv"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/helpers"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	usersService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/users_service"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/etc"
	"github.com/google/uuid"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

// @Router /v1/followings/follow [post]
// @Summary Follow a user
// @Description API for follow a user
// @Tags followings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.FollowRequest true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) Follow(c *gin.Context) {
	var (
		req models.FollowRequest
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	if req.FollowerId == req.FollowingId {
		c.JSON(200, models.Response{
			Code:    config.ErrorBadRequest,
			Message: "You can't follow yourself",
			Data:    nil,
		})
		return
	}

	fmt.Println(req.FollowerId, req.FollowingId)

	respFollow, err := h.services.FollowingService().FollowingService(ctx).Follow(ctx, &usersService.FollowRequest{
		FollowerId:  req.FollowerId,
		FollowingId: req.FollowingId,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error follow a user") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "sussecsfully done",
		Data:    respFollow,
	})
}

// @Router /v1/followings/unfollow [post]
// @Summary UnFollow a user
// @Description API for unfollow a user
// @Tags followings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.UnFollowRequest true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) UnFollow(c *gin.Context) {
	var (
		req models.UnFollowRequest
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	if req.FollowerId == req.FollowingId {
		c.JSON(200, models.Response{
			Code:    config.ErrorBadRequest,
			Message: "You can't unfollow yourself",
			Data:    nil,
		})
		return
	}

	respFollow, err := h.services.FollowingService().FollowingService(ctx).UnFollow(ctx, &usersService.UnFollowRequest{
		FollowerId:  req.FollowerId,
		FollowingId: req.FollowingId,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error unfollow a user") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "sussecsfully done",
		Data:    respFollow,
	})
}

// @Router /v1/followings/list-followings [get]
// @Summary Get Followings List With Pagination
// @Description API for getting followings list with pagination and search
// @Tags followings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id query string true "user id"
// @Param page query int false "Page number {default: 1}"
// @Param limit query int false "Limit number {default: 10}"
// @Param search query string false "Search keyword (name, surname, username, bio, etc.)"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetListFollowings(c *gin.Context) {
	ID, err := uuid.Parse(c.Query("id"))
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid user id") {
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting limit from query params") {
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting page from query params") {
		return
	}

	req := &usersService.GetListFollowingRequest{
		Id:    ID.String(),
		Limit: int32(limit),
		Page:  int32(page),
	}

	// Extract optional search and filters
	if search := c.Query("search"); search != "" {
		req.Search = &search
	}

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	h.log.Info("Received get list followings request", logger.Any("req", req))

	resp, err := h.services.FollowingService().FollowingService(ctx).GetListFollowings(ctx, req)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting followings list with pagination") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Get list followings successfully",
		Data:    resp,
	})
}

// @Router /v1/followings/list-followers [get]
// @Summary Get Followers List With Pagination
// @Description API for getting followers list with pagination and search
// @Tags followings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id query string true "user id"
// @Param page query int false "Page number {default: 1}"
// @Param limit query int false "Limit number {default: 10}"
// @Param search query string false "Search keyword (name, surname, username, bio, etc.)"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetListFollowers(c *gin.Context) {
	ID, err := uuid.Parse(c.Query("id"))
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid user id") {
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting limit from query params") {
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting page from query params") {
		return
	}

	req := &usersService.GetListFollowerRequest{
		Id:    ID.String(),
		Limit: int32(limit),
		Page:  int32(page),
	}

	// Extract optional search and filters
	if search := c.Query("search"); search != "" {
		req.Search = &search
	}

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	h.log.Info("Received get list followers request", logger.Any("req", req))

	resp, err := h.services.FollowingService().FollowingService(ctx).GetListFollowers(ctx, req)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting followers list with pagination") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Get list followers successfully",
		Data:    resp,
	})
}
