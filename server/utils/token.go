package utils

import (
	"fmt"
	"os"
	"tablelink_project/server/model"
	"time"

	"github.com/golang-jwt/jwt"
)

type ContextKey string

const (
	UserCtxKey              ContextKey = "user_id"
	AccessTokenExpiredTime             = 4 * time.Hour
	RefreshTokenExpiredTime            = 30 * 24 * time.Hour
	user_issuer                        = "tablelink_user"
)

type UserClaim struct {
	UserID       uint `json:"user_id"`
	RefreshToken bool `json:"rt"`
	jwt.StandardClaims
}

type RefreshTokenClaim struct {
	UserID       uint `json:"user_id"`
	RefreshToken bool `json:"rt"`
	jwt.StandardClaims
}

func GenerateToken(user_id uint) (model.Token, error) {
	expiredAt := time.Now().Add(AccessTokenExpiredTime)
	accessToken, err := GenerateAccessToken(user_id, expiredAt)
	if err != nil {
		return model.Token{}, err
	}

	refreshExpiredAt := time.Now().Add(RefreshTokenExpiredTime)
	refreshToken, err := GenerateRefreshToken(user_id, refreshExpiredAt)
	if err != nil {
		return model.Token{}, err
	}

	return model.Token{
		UserID:       user_id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    expiredAt,
	}, nil
}

func GenerateAccessToken(user_id uint, expiredAt time.Time) (string, error) {
	claims := createClaims(user_id, expiredAt)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func createClaims(user_id uint, expiredAt time.Time) jwt.Claims {
	return &UserClaim{
		UserID:       user_id,
		RefreshToken: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    user_issuer,
		},
	}
}

func GenerateRefreshToken(user_id uint, expiredAt time.Time) (string, error) {
	claims := createRefreshClaims(user_id, expiredAt)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func createRefreshClaims(user_id uint, expiredAt time.Time) jwt.Claims {
	return &RefreshTokenClaim{
		UserID:       user_id,
		RefreshToken: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    user_issuer,
		},
	}
}

func ValidateToken(tokenString string) (*UserClaim, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(*UserClaim)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return claim, nil
}

func ValidateRefreshToken(tokenString string) (*RefreshTokenClaim, error) {
	claim, token, err := ParseRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if !claim.RefreshToken {
		return nil, fmt.Errorf("invalid token")
	}

	return claim, nil
}

func ParseRefreshToken(tokenString string) (*RefreshTokenClaim, *jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, nil, err
	}

	claim, ok := token.Claims.(*RefreshTokenClaim)
	if !ok {
		return nil, nil, fmt.Errorf("invalid token")
	}

	return claim, token, nil
}
