package config

import (
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	AuthService
	UserService
	DB
	HTTPPort   string `validate:"required"`
	LogVerbose bool
}

type UserService struct {
	VerifyRegCode bool
}
type AuthService struct {
	JWTSecret string `validate:"required"`
}

type DB struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,numeric"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Database string `validate:"required"`
	SSL      string `validate:"required,oneof=enable disable"`
}

func New() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	c := Config{
		LogVerbose: v.GetBool("LOG_VERBOSE"),
		HTTPPort:   v.GetString("HTTP_PORT"),
		DB: DB{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetString("DB_PORT"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASS"),
			Database: v.GetString("DB_SCHEMA"),
			SSL:      v.GetString("DB_SSL"),
		},
		AuthService: AuthService{
			JWTSecret: v.GetString("AUTH_JWT_SECRET"),
		},
		UserService: UserService{
			VerifyRegCode: v.GetBool("USER_VERIFY_REG_CODE"),
		},
	}

	if err := validator.New().Struct(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
