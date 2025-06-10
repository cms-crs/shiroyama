package service

import (
	"authservice/src/dto"
	"authservice/src/model"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user model.User) (uint, error)
	GetUser(ctx context.Context, email string) (*model.User, error)
}

type AuthService struct {
	logger         *slog.Logger
	authRepository AuthRepository
}

func NewAuthService(authRepository AuthRepository, logger *slog.Logger) *AuthService {
	return &AuthService{
		logger:         logger,
		authRepository: authRepository,
	}
}

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) (*dto.RegisterResponse, error) {
	const op = "auth.service.Register"

	log := s.logger.With(
		slog.String("op", op),
	)

	email := request.Email
	password := request.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password")
		return nil, err
	}

	user := model.User{
		Email:    email,
		Password: hash,
	}

	id, err := s.authRepository.CreateUser(ctx, user)
	if err != nil {
		log.Error("Error creating user", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	return &dto.RegisterResponse{
		UserId: id,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	const op = "auth.service.Login"

	log := s.logger.With(
		slog.String("op", op),
	)

	email := request.Email
	password := request.Password

	user, err := s.authRepository.GetUser(ctx, email)
	if err != nil {
		log.Error("Error getting user", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		log.Error("Invalid Credentials", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	err = godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))

	return &dto.LoginResponse{
		Token: tokenString,
	}, err
}
