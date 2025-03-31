package v1

import (
	"bytes"
	"image"
	"strconv"
	"strings"

	// "github.com/aws/aws-sdk-go/aws"
	// awscredentials "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	rediscache "github.com/golanguzb70/redis-cache"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/logger"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/services"
)

type handlerV1 struct {
	log      logger.Logger
	cfg      *config.Config
	Redis    *redis.Client
	services services.ServiceManager
	cache    rediscache.RedisCache
}

type HandlerV1Options struct {
	Log      logger.Logger
	Cfg      *config.Config
	Redis    *redis.Client
	Services services.ServiceManager
	Cache    rediscache.RedisCache
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		log:      options.Log,
		cfg:      options.Cfg,
		services: options.Services,
		cache:    options.Cache,
		Redis:    options.Redis,
	}
}

func (h *handlerV1) Services() services.ServiceManager {
	return h.services
}

func (h *handlerV1) Log() logger.Logger {
	return h.log
}

func (h *handlerV1) Config() *config.Config {
	return h.cfg
}

func (h *handlerV1) Cache() rediscache.RedisCache {
	return h.cache
}

func (h *handlerV1) RedisClient() *redis.Client {
	return h.Redis
}

// ParsePageQueryParam ...
func ParsePageQueryParam(c *gin.Context) (uint64, error) {
	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		return 0, err
	}
	if page == 0 {
		return 1, nil
	}
	return page, nil
}

// ParseLimitQueryParam ...
func ParseLimitQueryParam(c *gin.Context) (uint64, error) {
	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 32)
	if err != nil {
		return 0, err
	}
	if limit == 0 {
		return 10, nil
	}
	return limit, nil
}

func generateFileName() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

func IsImageFile(data []byte) bool {
	_, _, err := image.Decode(bytes.NewBuffer(data))
	return err == nil
}

// func uploadToS3(cfg config.Config, filename string, data []byte, bucketName string) (string, error) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(cfg.CdnRegion),
// 		Credentials: awscredentials.NewStaticCredentials(
// 			cfg.CdnAccessKeyID,
// 			cfg.CdnSecretAccessKey,
// 			""),
// 		Endpoint: &cfg.CdnEndpoint,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	svc := s3.New(sess)

// 	_, err = svc.PutObject(&s3.PutObjectInput{
// 		Body:   bytes.NewReader(data),
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(filename),
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	imageURL := filename
// 	return imageURL, nil
// }
