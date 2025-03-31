package static

import (
	"time"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
)

var (
	AccessTokenSecret  = []byte(config.Load().AccessTokenSecret)
	RefreshTokenSecret = []byte(config.Load().RefreshTokenSecret)

	RefreshTokenKeyFormat = "%s user refresh token"
	UserCacheFormat       = "%s user cache"
)

const (
	AccessTokenExpireTime  = time.Hour * 24
	RefreshTokenExpireTime = time.Hour * 24 * 7

	CachingExpireTime = time.Hour * 24 * 7
)

const (
	TempleteBaseURL = "http://TaklifnomaVIP.uz/taklifnoma/"
)

var (
	MiniOEndpoint         = "minio.taklifnomavip.uz"
	MiniOAccessKey        = "minioadmin"
	MiniOSecretKey        = "minioadmin"
	MiniOBucket    string = "my-bucket"
	MiniORegion           = "us-east-1" // MinIO uses a default region
	MiniOProtocol         = true
)
