package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/v-starostin/goph-keeper/internal/model"
)

type Storage interface {
	AddUser(ctx context.Context, login, password string) error
	GetUser(ctx context.Context, login, password string) (*model.User, error)
	GetTokenByUserID(ctx context.Context, userID int32) (string, error)
	SaveToken(ctx context.Context, userID int32, token string) error
}

type Auth struct {
	storage Storage
	secret  []byte
}

func NewAuth(s Storage, sc []byte) *Auth {
	return &Auth{
		storage: s,
		secret:  sc,
	}
}

func (a *Auth) Register(ctx context.Context, username, password string) error {
	return a.storage.AddUser(ctx, username, password)
}

func (a *Auth) Authenticate(ctx context.Context, username, password string) (string, string, error) {
	user, err := a.storage.GetUser(ctx, username, password)
	if err != nil {
		return "", "", err
	}
	accessToken, err := a.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := a.generateRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = a.storage.SaveToken(ctx, user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *Auth) Refresh(ctx context.Context, accessToken, refreshToken string) (string, string, error) {
	token, err := jwt.ParseString(accessToken, jwt.WithVerify(jwa.HS256, a.secret))
	if err != nil {
		return "", "", err
	}
	userID, err := strconv.Atoi(token.Subject())
	if err != nil {
		return "", "", err
	}
	storedRefresh, err := a.storage.GetTokenByUserID(ctx, int32(userID))
	if err != nil {
		return "", "", err
	}
	if storedRefresh != refreshToken {
		return "", "", fmt.Errorf("wrong refresh token")
	}
	newAccessToken, err := a.generateAccessToken(int32(userID))
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := a.generateRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = a.storage.SaveToken(ctx, int32(userID), newRefreshToken)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (a *Auth) generateAccessToken(id int32) (string, error) {
	token := jwt.New()
	token.Set(jwt.SubjectKey, id)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().Add(10*time.Minute).Unix())
	tokenString, err := jwt.Sign(token, jwa.HS256, a.secret)
	if err != nil {
		return "", err
	}
	return string(tokenString), nil
}

func (a *Auth) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
