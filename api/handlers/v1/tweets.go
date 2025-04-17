package v1

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/helpers"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	tweetsService "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/genproto/tweets_service"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/etc"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/logger"
	static "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/statics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// @Router /v1/tweets [post]
// @Summary Tweet Create
// @Description API for creating tweet
// @Tags tweets
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.CreateTweet true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) CreateTweet(c *gin.Context) {
	var (
		req models.CreateTweet
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	respCreateTweet, err := h.services.TweetService().TweetService(ctx).CreateTweet(ctx, &tweetsService.CreateTweetRequest{
		Tweet: &tweetsService.Tweet{
			Id:      "",
			Content: req.Content,
			UserId:  req.UserId,
		},
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error creating tweet") {
		return
	}

	// add tweet to redis
	tweetData := models.Tweet{
		Id:      respCreateTweet.Id,
		Content: req.Content,
		UserId:  req.UserId,
	}
	err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, respCreateTweet.Id), etc.MarshalJSON(tweetData), static.CachingExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
		return
	}
	fmt.Println("Tweet added to redis successfully")

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet created sussecsfully",
		Data:    respCreateTweet,
	})
}

// @Router /v1/tweets/{id} [get]
// @Summary Get Tweet By ID
// @Description API for getting Tweet by id
// @Tags tweets
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetTweetByID(c *gin.Context) {
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
	tweetData, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetCacheFormat, id)).Result()
	if err == redis.Nil {
		resp, err := h.services.TweetService().TweetService(ctx).GetTweet(ctx, &tweetsService.GetTweetRequest{
			Id: id,
		})
		if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting tweet by id") {
			return
		}

		err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, id), etc.MarshalJSON(resp), static.CachingExpireTime).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
			return
		}
		fmt.Println("tweet data added to redis successfully")

		c.JSON(200, models.Response{
			Code:    config.StatusSuccess,
			Message: "Tweet got sussecsfully",
			Data:    resp,
		})
		return
	} else if err != nil {
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting tweet from redis") {
			return
		}
	}

	var resp models.Tweet
	err = json.Unmarshal([]byte(tweetData), &resp)
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshaling tweet data") {
		return
	}
	fmt.Println("got tweet data from redis successfully")

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet got sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/tweets/list [get]
// @Summary Get Tweet List With Pagination
// @Description API for getting Tweet list with pagination and filters
// @Tags tweets
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "Page number {default: 1}"
// @Param limit query int false "Limit number {default: 10}"
// @Param search query string false "Search keyword (role, login, fullname, phone, etc.)"
// @Param sort_by query string false "Sort field"
// @Param order query string false "Sorting order (asc or desc)"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetListTweets(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting limit from query params") {
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error getting page from query params") {
		return
	}

	req := &tweetsService.GetListTweetRequest{
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

	h.log.Info("Received get list tweet request", logger.Any("req", req))

	resp, err := h.services.TweetService().TweetService(ctx).GetListTweet(ctx, req)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting tweets list with pagination") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweets list retrieved successfully",
		Data:    resp,
	})
}

// @Router /v1/tweets [put]
// @Summary Tweet Update
// @Description API for updating Tweet
// @Tags tweets
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.UpdateTweet true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) UpdateTweet(c *gin.Context) {
	var (
		req models.UpdateTweet
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	resp, err := h.services.TweetService().TweetService(ctx).UpdateTweet(ctx, &tweetsService.UpdateTweetRequest{
		Tweet: &tweetsService.Tweet{
			Id:      req.Id,
			Content: req.Content,
			UserId:  req.UserId,
		},
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error updating tweet") {
		return
	}

	// getting tweet from redis for getting medias
	tweetData, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetCacheFormat, req.Id)).Result()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting tweet from redis") {
		return
	}
	var tweet models.Tweet
	err = json.Unmarshal([]byte(tweetData), &tweet)
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshaling tweet data") {
		return
	}
	tweet.Id = req.Id
	tweet.Content = req.Content
	tweet.UserId = req.UserId

	// add tweet to redis
	err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, req.Id), etc.MarshalJSON(tweet), static.CachingExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet updated sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/tweets/{id} [delete]
// @Summary Delete Tweet
// @Description API for deleting Tweet
// @Tags tweets
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) DeleteTweet(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	resp, err := h.services.TweetService().TweetService(ctx).DeleteTweet(ctx, &tweetsService.DeleteTweetRequest{
		Id: id,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error deleting tweet") {
		return
	}

	// delete tweet from redis
	err = h.Redis.Del(ctx, fmt.Sprintf(static.TweetCacheFormat, id)).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while deleting tweet from redis") {
		return
	}
	fmt.Println("deleted from redis successfully")

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet deleted sussecsfully",
		Data:    resp,
	})
}
