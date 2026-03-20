package usecase

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"baby-prep-quiz/domain"
)

var (
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Claims struct {
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AuthUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (u *AuthUsecase) SignUp(name, email, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.Create(name, email, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}
	return user, nil
}

func (u *AuthUsecase) Login(email, password string) (*domain.User, string, error) {
	user, passwordHash, err := u.userRepo.FindByEmail(email)
	if err == sql.ErrNoRows {
		return nil, "", ErrInvalidCredentials
	}
	if err != nil {
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := u.generateToken(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (u *AuthUsecase) ParseToken(tokenStr string) (*Claims, error) {
	c := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, c, func(t *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return c, nil
}

func (u *AuthUsecase) generateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString([]byte(u.jwtSecret))
}
