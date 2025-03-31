package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	v1 "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/handlers/v1"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	jwt "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/jwt"
)

type casbinPermission struct {
	enforcer *casbin.Enforcer
}

func (c *casbinPermission) CheckPermission(ctx *gin.Context, role string) (bool, error) {
	subject := role
	object := ctx.FullPath()
	action := ctx.Request.Method

	allow, err := c.enforcer.Enforce(subject, object, action)
	if err != nil {
		return false, err
	}

	return allow, nil
}

func Authentication(handlerOptions *v1.HandlerV1Options, cfg *config.Config, enf *casbin.Enforcer) func(ctx *gin.Context) {
	casbinHandler := &casbinPermission{
		enforcer: enf,
	}

	return func(ctx *gin.Context) {
		isOpen := config.OpenApis[ctx.FullPath()]
		if isOpen {
			ctx.Next()
			return
		}

		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(401, models.Response{
				Code:    config.ErrorCodeUnauthorized,
				Message: "Unauthorized request",
			})
			ctx.Abort()
			return
		}

		if len(strings.Split(token, " ")) > 1 {
			token = strings.Split(token, " ")[1]
		}

		tokenInfo, err := jwt.JWTExtract(token, cfg.AccessTokenSecret)
		if err != nil {
			ctx.JSON(401, models.Response{
				Code:    config.ErrorCodeUnauthorized,
				Message: err.Error(),
			})
			ctx.Abort()
			return
		}

		// for key, value := range tokenInfo {
		// 	ctx.Request.Header.Set(key, value)
		// 	key1 := key
		// 	value1 := value
		// 	ctx.Set(key1, value1)
		// }

		if tokenInfo["role"] == "" {
			log.Println("Invalid access token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Message: "Invalid access token",
			})
			return
		}

		ok, err := casbinHandler.CheckPermission(ctx, tokenInfo["role"])
		if err != nil {
			log.Println("Error while checking permission")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Message: "Error while checking permission",
				Data:    err.Error(),
			})
			return
		}

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.Response{
				Code:    config.ErrorCodeForbidden,
				Message: "You dont have right access",
			})
			return
		}

		ctx.Set("claims", tokenInfo)
		ctx.Next()
	}
}
