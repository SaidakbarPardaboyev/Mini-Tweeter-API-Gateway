package v1

import (
	"encoding/json"
	"fmt"

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

// @Router /v1/tweet-media [post]
// @Summary Tweet Media Create
// @Description API for creating tweet Media
// @Tags tweet-media
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.CreateTweetMedia true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) CreateTweetMedia(c *gin.Context) {
	var (
		req models.CreateTweetMedia
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	respCreateTweetMedia, err := h.services.TweetService().TweetMediaService(ctx).CreateTweetMedia(ctx, &tweetsService.CreateTweetMediaRequest{
		TweetMedia: &tweetsService.TweetMedia{
			TweetId:   req.TweetId,
			MediaType: req.MediaType,
			Url:       req.Url,
		},
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error creating tweet media") {
		return
	}

	// add tweet media to redis
	{
		// add tweet media to redis
		tweetMediaData := models.TweetMedia{
			Id:        respCreateTweetMedia.Id,
			TweetId:   req.TweetId,
			MediaType: req.MediaType,
			Url:       req.Url,
		}
		err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetMediaCacheFormat, respCreateTweetMedia.Id), etc.MarshalJSON(tweetMediaData), static.CachingExpireTime).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet media to redis") {
			return
		}
		fmt.Println("Tweet media added to redis successfully")
	}

	// add media to tweet redis data
	{
		// getting tweet from redis
		tweetDataStr, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetCacheFormat, req.TweetId)).Result()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting tweet from redis") {
			return
		}

		// unmarshal tweet data
		var tweet models.Tweet
		err = json.Unmarshal([]byte(tweetDataStr), &tweet)
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshalling tweet data") {
			return
		}

		tweet.Medias = append(tweet.Medias, &models.TweetMedia{
			Id:        respCreateTweetMedia.Id,
			TweetId:   req.TweetId,
			MediaType: req.MediaType,
			Url:       req.Url,
		})
		err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, req.TweetId), etc.MarshalJSON(tweet), static.CachingExpireTime).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
			return
		}
		fmt.Println("Tweet added to redis successfully")
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet media created sussecsfully",
		Data:    respCreateTweetMedia,
	})
}

// @Router /v1/tweet-media/{id} [get]
// @Summary Get Tweet Media By ID
// @Description API for getting Tweet Media by id
// @Tags tweet-media
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetTweetMediaByID(c *gin.Context) {
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
	tweetMediaData, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetMediaCacheFormat, id)).Result()
	if err == redis.Nil {
		resp, err := h.services.TweetService().TweetMediaService(ctx).GetTweetMedia(ctx, &tweetsService.GetTweetMediaRequest{
			Id: id,
		})
		if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting tweet media by id") {
			return
		}

		err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetMediaCacheFormat, id), etc.MarshalJSON(resp), static.CachingExpireTime).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet media to redis") {
			return
		}
		fmt.Println("tweet media data added to redis successfully")

		c.JSON(200, models.Response{
			Code:    config.StatusSuccess,
			Message: "Tweet media got sussecsfully",
			Data:    resp,
		})
		return
	} else if err != nil {
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while getting tweet media from redis") {
			return
		}
	}

	var resp models.TweetMedia
	err = json.Unmarshal([]byte(tweetMediaData), &resp)
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshaling tweet media data") {
		return
	}
	fmt.Println("got tweet media data from redis successfully")

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet media got sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/tweet-media/list [get]
// @Summary Get Tweet Media List With Pagination
// @Description API for getting Tweet Media list with pagination and filters
// @Tags tweet-media
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param tweet_id query string true "Tweet ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) GetListTweetMedia(c *gin.Context) {
	tweetID := c.DefaultQuery("tweet_id", "")
	if tweetID == "" {
		if helpers.HandleBadRequestErrWithMessage(c, h.log, fmt.Errorf("tweet_id is required"), "Invalid request body") {
			return
		}
	}

	req := &tweetsService.GetListTweetMediaRequest{
		TweetId: tweetID,
	}

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	h.log.Info("Received get list tweet media request", logger.Any("req", req))

	resp, err := h.services.TweetService().TweetMediaService(ctx).GetListTweetMedia(ctx, req)
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error getting tweets medias list with pagination") {
		return
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweets medias list retrieved successfully",
		Data:    resp,
	})
}

// @Router /v1/tweet-media [put]
// @Summary Tweet Media Update
// @Description API for updating Tweet Media
// @Tags tweet-media
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.UpdateTweetMedia true "request"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) UpdateTweetMedia(c *gin.Context) {
	var (
		req models.UpdateTweetMedia
	)

	err := c.ShouldBindJSON(&req)
	if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "Invalid request body") {
		return
	}
	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	resp, err := h.services.TweetService().TweetMediaService(ctx).UpdateTweetMedia(ctx, &tweetsService.UpdateTweetMediaRequest{
		TweetMedia: &tweetsService.TweetMedia{
			Id:        req.Id,
			TweetId:   req.TweetId,
			MediaType: req.MediaType,
			Url:       req.Url,
		},
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error updating tweet media") {
		return
	}

	// add tweet media to redis
	err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetMediaCacheFormat, req.Id), etc.MarshalJSON(req), static.CachingExpireTime).Err()
	if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet media to redis") {
		return
	}

	// update image tweet media from tweet's medias redis
	{
		// getting tweet from redis
		tweetData, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetCacheFormat, req.TweetId)).Result()
		if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "tweet not found") {
			return
		}

		var tweet models.Tweet
		err = json.Unmarshal([]byte(tweetData), &tweet)
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshalling tweet from redis") {
			return
		}

		var found bool
		for i, media := range tweet.Medias {
			if media.Id == req.Id {
				tweet.Medias[i] = &models.TweetMedia{
					Id:        req.Id,
					TweetId:   req.TweetId,
					MediaType: req.MediaType,
					Url:       req.Url,
				}
				found = true
				break
			}
		}

		if found {
			// add tweet to redis
			err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, req.TweetId), etc.MarshalJSON(tweet), static.CachingExpireTime).Err()
			if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
				return
			}
		}
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet media updated sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/tweet-media/{tweet_id}/{id} [delete]
// @Summary Delete Tweet Media
// @Description API for deleting Tweet Media
// @Tags tweet-media
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param tweet_id path string true "tweet id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) DeleteTweetMedia(c *gin.Context) {
	id := c.Param("id")
	tweetID := c.Param("tweet_id")
	fmt.Println(id, "TwetID", tweetID)

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	resp, err := h.services.TweetService().TweetMediaService(ctx).DeleteTweetMedia(ctx, &tweetsService.DeleteTweetMediaRequest{
		Id: id,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error deleting tweet media") {
		return
	}

	{
		// delete tweet from redis
		err = h.Redis.Del(ctx, fmt.Sprintf(static.TweetMediaCacheFormat, id)).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while deleting tweet media from redis") {
			return
		}
		fmt.Println("deleted from redis successfully")
	}

	// delete tweet media from tweet's medias redis
	{
		// getting tweet from redis
		tweetData, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetCacheFormat, tweetID)).Result()
		if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "tweet not found") {
			return
		}

		var tweet models.Tweet
		err = json.Unmarshal([]byte(tweetData), &tweet)
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshalling tweet from redis") {
			return
		}

		var found bool
		for i, media := range tweet.Medias {
			if media.Id == id {
				tweet.Medias = append(tweet.Medias[:i], tweet.Medias[i+1:]...)
				found = true
				break
			}
		}

		if found {
			// add tweet to redis
			err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, tweetID), etc.MarshalJSON(tweet), static.CachingExpireTime).Err()
			if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
				return
			}
		}
	}

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet media deleted sussecsfully",
		Data:    resp,
	})
}

// @Router /v1/tweet-media/tweet/{id} [delete]
// @Summary Delete Tweet Media With Tweet ID
// @Description API for deleting Tweet Media with tweet id
// @Tags tweet-media
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path string true "tweet id"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 403 {object} models.Response "Forbidden"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Failure default {object} models.Response
func (h *handlerV1) DeleteTweetMediaWithTweetID(c *gin.Context) {
	tweetID := c.Param("id")

	ctx, cancel := etc.NewTimoutContext(c.Request.Context())
	defer cancel()

	resp, err := h.services.TweetService().TweetMediaService(ctx).DeleteTweetMediaWithTweetID(ctx, &tweetsService.DeleteTweetMediaWithTweetIDRequest{
		TweetId: tweetID,
	})
	if helpers.HandleGrpcErrWithMessage(c, h.log, err, "Error deleting tweet medias with tweet ID") {
		return
	}

	// make tweet's medias equel to empty array
	{
		// getting tweet from redis
		tweetData, err := h.Redis.Get(ctx, fmt.Sprintf(static.TweetCacheFormat, tweetID)).Result()
		if helpers.HandleBadRequestErrWithMessage(c, h.log, err, "tweet not found") {
			return
		}

		var tweet models.Tweet
		err = json.Unmarshal([]byte(tweetData), &tweet)
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while unmarshalling tweet from redis") {
			return
		}
		tweet.Medias = []*models.TweetMedia{}

		// add tweet to redis
		err = h.Redis.Set(ctx, fmt.Sprintf(static.TweetCacheFormat, tweetID), etc.MarshalJSON(tweet), static.CachingExpireTime).Err()
		if helpers.HandleInternalWithMessage(c, h.log, err, "Error while adding tweet to redis") {
			return
		}
	}
	fmt.Println("deleted medias from tweet's medias from redis successfully")

	c.JSON(200, models.Response{
		Code:    config.StatusSuccess,
		Message: "Tweet medias deleted sussecsfully",
		Data:    resp,
	})
}
