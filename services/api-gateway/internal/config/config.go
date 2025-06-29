package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port         int    `json:"port"`
	Environment  string `json:"environment"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`

	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`

	JWTSecret            string `json:"jwt_secret"`
	AccessTokenDuration  int    `json:"access_token_duration"`
	RefreshTokenDuration int    `json:"refresh_token_duration"`

	Services ServiceConfig `json:"services"`

	RateLimitEnabled bool `json:"rate_limit_enabled"`
	RateLimitRPS     int  `json:"rate_limit_rps"`

	LogLevel  string `json:"log_level"`
	LogFormat string `json:"log_format"`
}

type ServiceConfig struct {
	AuthService     string `json:"auth_service"`
	UserService     string `json:"user_service"`
	TeamService     string `json:"team_service"`
	BoardService    string `json:"board_service"`
	TaskService     string `json:"task_service"`
	CommentService  string `json:"comment_service"`
	ActivityService string `json:"activity_service"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnvInt("PORT", 8080),
		Environment:  getEnv("ENVIRONMENT", "development"),
		ReadTimeout:  getEnvInt("READ_TIMEOUT", 30),
		WriteTimeout: getEnvInt("WRITE_TIMEOUT", 30),
		IdleTimeout:  getEnvInt("IDLE_TIMEOUT", 120),

		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ","),
		AllowedMethods: strings.Split(getEnv("ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"), ","),
		AllowedHeaders: strings.Split(getEnv("ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"), ","),

		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key"),
		AccessTokenDuration:  getEnvInt("ACCESS_TOKEN_DURATION", 3600),    // 1 hour
		RefreshTokenDuration: getEnvInt("REFRESH_TOKEN_DURATION", 604800), // 7 days

		Services: ServiceConfig{
			AuthService:     getEnv("AUTH_SERVICE_URL", "localhost:50001"),
			UserService:     getEnv("USER_SERVICE_URL", "localhost:50002"),
			TeamService:     getEnv("TEAM_SERVICE_URL", "localhost:50003"),
			BoardService:    getEnv("BOARD_SERVICE_URL", "localhost:50004"),
			TaskService:     getEnv("TASK_SERVICE_URL", "localhost:50005"),
			CommentService:  getEnv("COMMENT_SERVICE_URL", "localhost:50006"),
			ActivityService: getEnv("ACTIVITY_SERVICE_URL", "localhost:50007"),
		},

		RateLimitEnabled: getEnvBool("RATE_LIMIT_ENABLED", false),
		RateLimitRPS:     getEnvInt("RATE_LIMIT_RPS", 100),

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
