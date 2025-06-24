package service

import (
	"authservice/src/config"
	"authservice/src/dto"
	"authservice/src/model"
	"context"
	"github.com/cms-crs/protos/gen/go/user_service"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user model.User) (uint, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateRefreshToken(ctx context.Context, user *model.User, token string) error
	GetRefreshToken(ctx context.Context, userId string) (string, error)
}

type AuthService struct {
	logger         *slog.Logger
	authRepository AuthRepository
	config         *config.Config
	userClient     userv1.UserServiceClient
}

type TokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(authRepository AuthRepository, logger *slog.Logger, config *config.Config, userConn *grpc.ClientConn) *AuthService {

	return &AuthService{
		logger:         logger,
		authRepository: authRepository,
		config:         config,
		userClient:     userv1.NewUserServiceClient(userConn),
	}
}

func (s *AuthService) Register(ctx context.Context, request dto.RegisterRequest) (*dto.RegisterResponse, error) {
	const op = "auth.service.Register"

	log := s.logger.With(
		slog.String("op", op),
	)

	email := request.Email
	password := request.Password
	username := request.Username

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password")
		return nil, err
	}

	userClientUser, err := s.userClient.CreateUser(ctx, &userv1.CreateUserRequest{
		Email:    email,
		Username: username,
	})
	if err != nil {
		log.Error("Error creating user")
		return nil, err
	}

	user := model.User{
		Email:    email,
		Password: hash,
		UserID:   userClientUser.GetId(),
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

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return at.SignedString([]byte(secret))
}

func (s *AuthService) generateRefreshToken(user *model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(s.config.JWT.RefreshTokenTTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return rt.SignedString([]byte(secret))
}

func (s *AuthService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	const op = "auth.service.Login"

	log := s.logger.With(
		slog.String("op", op),
	)

	email := request.Email
	password := request.Password

	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("Error getting user", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		log.Error("Invalid Credentials", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	// access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// refresh token
	refreshToken, err := s.generateRefreshToken(user)

	user.RefreshToken = refreshToken

	err = s.authRepository.UpdateRefreshToken(ctx, user, refreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

func (s *AuthService) ValidateToken(_ context.Context, accessToken string) (bool, error) {
	const op = "auth.service.ValidateToken"
	log := s.logger.With(slog.String("op", op))

	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Error("Error parsing token")
		return false, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
			log.Error("Token is expired")
			return false, nil
		}
	}

	return true, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	const op = "auth.service.RefreshToken"
	log := s.logger.With(slog.String("op", op))

	isValid, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		log.Debug("validate token failed", "error", err)
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	if !isValid {
		log.Debug("validate token failed", "error", err)
		return "", status.Error(codes.PermissionDenied, "invalid refresh token")
	}

	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Error("Error parsing token")
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		log.Debug("parse refresh token failed", "error", err)
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	dbRefreshToken, err := s.authRepository.GetRefreshToken(ctx, claims.Subject)
	if err != nil {
		log.Debug("get refresh token failed", "error", err)
		return "", status.Error(codes.InvalidArgument, err.Error())
	}

	if dbRefreshToken != refreshToken {
		log.Debug("refresh token failed", "error", err)
		return "", status.Error(codes.PermissionDenied, "invalid refresh token")
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 0)
	user, err := s.authRepository.GetUserByID(ctx, uint(id))

	return s.generateAccessToken(user)
}
