package config

import (
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	AuthService
	UserService
	RequestService
	DB
	GuardUI
	S3
	SMTP
	HTTPPort string `validate:"required"`
	LogLevel string
}

type UserService struct {
	VerifyRegCode bool
	MembersLimit  int
}

type AuthService struct {
	JWTSecret string `validate:"required"`
}

type RequestService struct {
	S3SpaceName string `validate:"required"`
	CDNHost     string `validate:"required"`
}

type GuardUI struct {
	APIHost    string `validate:"required"`
	PageURI    string `validate:"required"`
	PagerLimit int    `validate:"required"`
}

type DB struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,numeric"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Database string `validate:"required"`
	SSL      string `validate:"required,oneof=enable disable require"`
}

type S3 struct {
	Region   string `validate:"required"`
	Key      string `validate:"required"`
	Secret   string `validate:"required"`
	Endpoint string `validate:"required"`
}

type SMTP struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
	From string `validate:"required"`
	Pass string `validate:"required"`
}

func New() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	c := Config{
		LogLevel: v.GetString("LOG_LEVEL"),
		HTTPPort: v.GetString("HTTP_PORT"),
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
			MembersLimit:  v.GetInt("FAMILY_MEMBERS_LIMIT"),
		},
		RequestService: RequestService{
			S3SpaceName: v.GetString("S3_SPACE_NAME"),
			CDNHost:     v.GetString("CDN_HOST"),
		},
		GuardUI: GuardUI{
			APIHost:    v.GetString("UI_GUARD_API_HOST"),
			PageURI:    v.GetString("UI_GUARD_PAGE_URI"),
			PagerLimit: v.GetInt("UI_GUARD_PAGER_LIMIT"),
		},
		S3: S3{
			Region:   v.GetString("S3_REGION"),
			Key:      v.GetString("S3_KEY"),
			Secret:   v.GetString("S3_SECRET"),
			Endpoint: v.GetString("S3_ENDPOINT"),
		},
		SMTP: SMTP{
			Host: v.GetString("SMTP_HOST"),
			Port: v.GetString("SMTP_PORT"),
			From: v.GetString("SMTP_FROM"),
			Pass: v.GetString("SMTP_PASS"),
		},
	}

	if err := validator.New().Struct(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
