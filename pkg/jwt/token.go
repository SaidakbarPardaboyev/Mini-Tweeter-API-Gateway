package etc

import (
	"fmt"
	"time"

	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/api/models"
	"github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/config"
	static "github.com/SaidakbarPardaboyev/Mini-Tweeter-API-Gateway/pkg/statics"
	"github.com/spf13/cast"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       user.Id,
		"name":          user.Name,
		"surname":       user.Surname,
		"bio":           user.Bio,
		"profile_image": user.ProfileImage,
		"username":      user.Username,
		"exp":           time.Now().Add(static.AccessTokenExpireTime).Unix(), // Expires in 15 min
		"iat":           time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(static.AccessTokenSecret)
}

func GenerateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       user.Id,
		"name":          user.Name,
		"surname":       user.Surname,
		"bio":           user.Bio,
		"profile_image": user.ProfileImage,
		"username":      user.Username,
		"exp":           time.Now().Add(static.RefreshTokenExpireTime).Unix(), // Expires in 7 days
		"iat":           time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(static.RefreshTokenSecret)
}

func ValidateToken(tokenString string, secretKey []byte) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, claims, nil
	}
	return nil, nil, fmt.Errorf("invalid token")
}

func RefreshToken(refreshToken string) (string, error) {
	_, claims, err := ValidateToken(refreshToken, static.RefreshTokenSecret)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	newAccessToken, err := GenerateAccessToken(&models.User{
		Id:           claims["user_id"].(string),
		Name:         claims["name"].(string),
		Surname:      claims["surname"].(string),
		Bio:          claims["bio"].(string),
		ProfileImage: claims["profile_image"].(string),
		Username:     claims["username"].(string),
	})
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func JWTExtract(tokenStr string, signKey string) (map[string]string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		err = fmt.Errorf("%s", config.ErrorCodeInvalidToken)
		return nil, err
	}

	res := map[string]string{}
	for k, v := range claims {
		res[k] = cast.ToString(v)
	}

	return res, nil
}
