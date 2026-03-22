package middlewares

import (
	"backend/src/config"
	"backend/src/models"
	"backend/src/services"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
    ErrExpiredToken       = errors.New("token has expired")
    ErrEmailInUse         = errors.New("email already in use")
)

type AuthService struct {
    userRepo         *services.UserStore
    refreshTokenRepo *models.UserRefreshToken
    jwtSecret        []byte
    accessTokenTTL   time.Duration
}

func NewAuthService(userRepo *services.UserStore, refreshTokenRepo *models.UserRefreshToken, cfg *config.ServerConfig, accessTokenTTL time.Duration) *AuthService {
    return &AuthService{
        userRepo:         userRepo,
        refreshTokenRepo: refreshTokenRepo,
        jwtSecret:        []byte(cfg.JWTSecret),
        accessTokenTTL:   accessTokenTTL,
    }
}

func (ac *AuthService) CreateToken(username, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
	jwt.MapClaims{
		"sub": id,
		"username": username,
		"exp":time.Now().Add(time.Hour * 3).Unix(),
	})

	tokenString, err := token.SignedString(ac.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ac *AuthService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		} 
			return ac.jwtSecret, nil 
	})

	
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

	return nil, ErrInvalidToken
}