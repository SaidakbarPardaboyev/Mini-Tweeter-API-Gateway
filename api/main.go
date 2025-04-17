package api

import (
	v1 "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/handlers/v1"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/middleware"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/logger"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/services"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	rediscache "github.com/golanguzb70/redis-cache"
	"github.com/redis/go-redis/v9"

	_ "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterOptions struct {
	Log      logger.Logger
	Cfg      *config.Config
	Redis    *redis.Client
	Services services.ServiceManager
	Cache    rediscache.RedisCache
	Enforcer *casbin.Enforcer
}

// @title MiniTweeter API Gateway v1
// @version 1.0
// @description This is a MiniTweeter API Gateway server.
// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @name Authorization
// @in header
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = false

	router.Use(CORSMiddleware())

	// Handle 404 Errors with CORS
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Route not found"})
	})

	options := &v1.HandlerV1Options{
		Log:      opt.Log,
		Cfg:      opt.Cfg,
		Services: opt.Services,
		Cache:    opt.Cache,
		Redis:    opt.Redis,
	}

	handlerV1 := v1.New(options)

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	v1Group := router.Group("/v1")
	v1Group.Use(middleware.Authentication(options, opt.Cfg, opt.Enforcer))

	auth := v1Group.Group("/auth")
	{
		auth.POST("/sign-up", handlerV1.SignUp)
		auth.POST("/login", handlerV1.Login)
		auth.POST("/logout", handlerV1.Logout)
		// auth.GET("/refresh-token", handlerV1.RefreshToken)
	}

	fileUpload := v1Group.Group("/file")
	{
		fileUpload.POST("/upload", handlerV1.UploadFile)
		fileUpload.DELETE("/delete", handlerV1.DeleteFile)
	}

	users := v1Group.Group("/users")
	{
		users.POST("", handlerV1.CreateUser)
		users.PUT("", handlerV1.UpdateUser)
		users.GET("/:id", handlerV1.GetUserByID)
		users.GET("/list", handlerV1.GetLislUsers)
		users.DELETE("/:id", handlerV1.DeleteUser)
	}

	followings := v1Group.Group("/followings")
	{
		followings.POST("/follow", handlerV1.Follow)
		followings.POST("/unfollow", handlerV1.UnFollow)
		followings.GET("/list-followings", handlerV1.GetListFollowings)
		followings.GET("/list-followers", handlerV1.GetListFollowers)
	}

	tweets := v1Group.Group("/tweets")
	{
		tweets.POST("", handlerV1.CreateTweet)
		tweets.GET("/list", handlerV1.GetListTweets)
		tweets.GET("/:id", handlerV1.GetTweetByID)
		tweets.PUT("", handlerV1.UpdateTweet)
		tweets.DELETE("/:id", handlerV1.DeleteTweet)
	}

	tweetMedia := v1Group.Group("/tweet-media")
	{
		tweetMedia.POST("", handlerV1.CreateTweetMedia)
		tweetMedia.GET("/:id", handlerV1.GetTweetMediaByID)
		tweetMedia.GET("/list", handlerV1.GetListTweetMedia)
		tweetMedia.PUT("", handlerV1.UpdateTweetMedia)
		tweetMedia.DELETE("/:tweet_id/:id", handlerV1.DeleteTweetMedia)
		tweetMedia.DELETE("/tweet/:id", handlerV1.DeleteTweetMediaWithTweetID)
	}

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, HEAD, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
