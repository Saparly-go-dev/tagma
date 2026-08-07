package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt     = "aslkdgoialskdfalsdflaksdjhf"
	tokenTTL = 7 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user tagma.User) (int, error) {
	if !isAllowedRole(user.Type) {
		return 0, errors.New("invalid user role")
	}
	password, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = password
	user.Status = true
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.authenticate(username, password)
	if err != nil {
		return "", err
	}

	signingKey, err := jwtSigningKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: user.Id,
		Role:   user.Type,
	})

	return token.SignedString(signingKey)
}

func (s *AuthService) GetUserType(username, password string) (string, error) {
	user, err := s.authenticate(username, password)
	if err != nil {
		return "", err
	}

	return user.Type, nil
}

func (s *AuthService) GetTradeAgentId(Id int) (int, error) {
	return s.repo.GetAgentId(Id)
}

func (s *AuthService) GetEkspeditorId(Id int) (int, error) {
	return s.repo.GetEkspeditorId(Id)
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	signingKey, err := jwtSigningKey()
	if err != nil {
		return 0, "", err
	}

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return signingKey, nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New("invalid token claims")
	}

	user, err := s.repo.GetUserByID(claims.UserId)
	if err != nil {
		return 0, "", err
	}
	if !user.Status {
		return 0, "", errors.New("user is disabled")
	}
	if !isAllowedRole(user.Type) {
		return 0, "", errors.New("invalid user role")
	}
	return claims.UserId, user.Type, nil
}

func (s *AuthService) GetTradeAgentIdFromEkspeditorId(Id int) (int, error) {
	return s.repo.GetTradeAgentIdFromEkspeditorId(Id)
}

func (s *AuthService) GetEkspeditoryIdFromAgentId(Id int) (int, error) {
	return s.repo.GetEkspeditoryIdFromAgentId(Id)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *AuthService) authenticate(username, password string) (tagma.User, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return tagma.User{}, err
	}
	if !user.Status {
		return tagma.User{}, errors.New("user is disabled")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return user, nil
	}

	if user.Password != generatePasswordHash(password) {
		return tagma.User{}, errors.New("invalid credentials")
	}

	newHash, err := hashPassword(password)
	if err != nil {
		return tagma.User{}, err
	}
	if err := s.repo.UpdatePassword(user.Id, newHash); err != nil {
		return tagma.User{}, err
	}
	user.Password = newHash
	return user, nil
}

func jwtSigningKey() ([]byte, error) {
	key := os.Getenv("JWT_SIGNING_KEY")
	if len(key) < 32 {
		return nil, errors.New("JWT_SIGNING_KEY must contain at least 32 characters")
	}
	return []byte(key), nil
}

func isAllowedRole(role string) bool {
	switch role {
	case "admin", "agent", "ekspeditor", "merchandiser", "viewer":
		return true
	default:
		return false
	}
}
