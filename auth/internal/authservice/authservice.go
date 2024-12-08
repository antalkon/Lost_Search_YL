package authservice

import (
	"auth/internal/userrepo"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	authServiceName = "auth_service"
)

type AuthRepoInterface interface {
	CreateUser(ctx context.Context, u userrepo.User) error
	GetUser(ctx context.Context, login string) (userrepo.User, error)
	UpdateUser(ctx context.Context, u userrepo.User) (userrepo.User, error)
}

type AuthServiceConfig struct {
	TokenExpiration int    `env:"TOKEN_EXPIRATION" env-default:"3600"`
	TokenSigningKey string `env:"TOKEN_SIGNING_KEY" env-default:"shikanokonokokoshitantan"`
}

type AuthService struct {
	Config   *AuthServiceConfig
	AuthRepo AuthRepoInterface
}

func NewAuthService(config *AuthServiceConfig, authRepo AuthRepoInterface) *AuthService {
	return &AuthService{
		Config:   config,
		AuthRepo: authRepo,
	}
}

var ErrUserAlreadyExists = fmt.Errorf("user already exists")

// создание пользователя
// возвращает jwt токен и ошибку
func (as *AuthService) CreateUser(ctx context.Context,
	login string, password string, data string) (string, error) {
	// проверка на существующего пользователя
	if u, _ := as.AuthRepo.GetUser(ctx, login); u != (userrepo.User{}) {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u := userrepo.User{
		Login:    login,
		Password: string(hashedPassword),
		Data:     data,
	}

	err = as.AuthRepo.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    authServiceName,
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(as.Config.TokenExpiration))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := claims.SignedString([]byte(as.Config.TokenSigningKey))

	return tokenString, err
}

// вход пользователя
func (as *AuthService) LoginUser(ctx context.Context, login string, password string) (string, error) {
	u, err := as.AuthRepo.GetUser(ctx, login)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    authServiceName,
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(as.Config.TokenExpiration))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := claims.SignedString([]byte(as.Config.TokenSigningKey))

	return tokenString, err
}

// проверка на валидность токена
func (as *AuthService) IsTokenValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.Config.TokenSigningKey), nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}
