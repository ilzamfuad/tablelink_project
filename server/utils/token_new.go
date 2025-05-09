package utils

import (
	"fmt"
	"strings"
	"tablelink_project/server/model"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	bearer                 = "Bearer "
	accessTokenExpireTime  = 4 * time.Hour
	refreshTokenExpireTime = 30 * 24 * time.Hour
	issuer                 = "Wartek-ID"
)

var (
	defaultSigningMethod = jwt.SigningMethodHS512
)

// CredentialTokenGenerator is responsible for managing credential tokens.
type CredentialTokenGenerator interface {
	// Generate generates token based on id.
	Generate(id uuid.UUID, name string) (*model.Token, error)
	// RenewToken renews the token using a valid refresh token.
	RenewToken(refreshToken string, name string) (*model.Token, error)
	// ParseToken parse the token
	ParseToken(accessToken string) (*UserClaims, *jwt.Token, error)
}

// UserClaims defines JWT claims used in Wartek-ID.
type UserClaims struct {
	UserID       uint `json:"user_id"`
	RefreshToken bool `json:"rt"`
	jwt.StandardClaims
}

// RefreshTokenClaims defines JWT refresh token claims used in Wartek-ID.
type RefreshTokenClaims struct {
	UserID       uint `json:"user_id"`
	RefreshToken bool `json:"rt"`
	jwt.StandardClaims
}

// JWT holds the job to generate token using JWT.
type JWT struct {
	secretKey string
}

// NewJWT creates an instance of JWT.
func NewJWT(secretKey string) *JWT {
	return &JWT{
		secretKey: secretKey,
	}
}

// Generate generates token based on id.
func (j *JWT) Generate(id uint, name string) (*model.Token, error) {
	accessTokenExpireAt := time.Now().Add(accessTokenExpireTime)
	accessToken, err := j.createAccessToken(accessTokenExpireAt, id, name)
	if err != nil {
		return nil, errors.Wrap(err, "[JWT-Generate] error creating access token")
	}

	refreshTokenExpireAt := time.Now().Add(refreshTokenExpireTime)
	refreshToken, err := j.createRefreshToken(refreshTokenExpireAt, id)
	if err != nil {
		return nil, errors.Wrap(err, "[JWT-Generate] error creating refresh token")
	}

	cred := &model.Token{
		UserID:       id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    accessTokenExpireAt,
	}
	return cred, nil
}

// IsRefreshTokenValid checks if refresh token is valid.
func (j *JWT) IsRefreshTokenValid(refreshToken string) (*RefreshTokenClaims, error) {
	claim, token, err := j.ParseRefreshToken(refreshToken)

	if err != nil {
		return nil, errors.Wrap(err, "[JWT-IsRefreshTokenValid] error parsing token")
	}
	if !token.Valid {
		return nil, errors.Wrap(err, "[JWT-IsRefreshTokenValid] token invalid")
	}
	if !claim.RefreshToken {
		return nil, errors.Wrap(err, "[JWT-IsRefreshTokenValid] token invalid")
	}
	return claim, nil
}

// RenewToken renews access and refresh token.
func (j *JWT) RenewToken(token string, name string) (*model.Token, error) {
	claim, err := j.IsRefreshTokenValid(token)
	if err != nil {
		return nil, errors.Wrap(err, "[JWT-RenewToken] error check token is valid")
	}
	return j.Generate(claim.UserID, name)
}

func (j *JWT) createAccessToken(expireAt time.Time, id uint, name string) (string, error) {
	return j.createToken(expireAt, id, name)
}

func (j *JWT) createRefreshToken(expireAt time.Time, id uint) (string, error) {
	claim := createRefreshTokenClaim(expireAt, id)
	token := jwt.NewWithClaims(defaultSigningMethod, claim)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) createToken(expireAt time.Time, id uint, name string) (string, error) {
	claim := createClaim(expireAt, id, name)
	token := jwt.NewWithClaims(defaultSigningMethod, claim)
	return token.SignedString([]byte(j.secretKey))
}

// ParseToken parse token
func (j *JWT) ParseToken(signedToken string) (*UserClaims, *jwt.Token, error) {
	if strings.Contains(signedToken, bearer) {
		signedToken = strings.TrimPrefix(signedToken, bearer)
	}
	token, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"]), "JWT-ParseToken")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "JWT-ParseToken")
	}

	claim, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, nil, errors.Wrap(errors.New("error claiming token"), "JWT-ParseToken")
	}
	return claim, token, nil
}

// ParseRefreshToken parse refresh token
func (j *JWT) ParseRefreshToken(signedToken string) (*RefreshTokenClaims, *jwt.Token, error) {
	if strings.Contains(signedToken, bearer) {
		signedToken = strings.TrimPrefix(signedToken, bearer)
	}
	token, err := jwt.ParseWithClaims(signedToken, &RefreshTokenClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"]), "JWT-ParseRefreshToken")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "JWT-ParseRefreshToken")
	}

	claim, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, nil, errors.Wrap(errors.New("error claiming token"), "JWT-ParseRefreshToken")
	}
	return claim, token, nil
}

func createClaim(expireAt time.Time, id uint, name string) jwt.Claims {
	return &UserClaims{
		UserID:       id,
		RefreshToken: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
}

func createRefreshTokenClaim(expireAt time.Time, id uint) jwt.Claims {
	return &RefreshTokenClaims{
		UserID:       id,
		RefreshToken: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
}
