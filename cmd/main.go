package main

import (
	"context"
	"errors"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"

	clientUsers "github.com/ivch/dynasty/common/clients/users"
	"github.com/ivch/dynasty/common/email"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/config"
	"github.com/ivch/dynasty/server"
	svcAuth "github.com/ivch/dynasty/server/handlers/auth"
	repoAuth "github.com/ivch/dynasty/server/handlers/auth/repo"
	transportAuth "github.com/ivch/dynasty/server/handlers/auth/transport"
	svcDict "github.com/ivch/dynasty/server/handlers/dictionaries"
	repoDict "github.com/ivch/dynasty/server/handlers/dictionaries/repo"
	transportDict "github.com/ivch/dynasty/server/handlers/dictionaries/transport"
	"github.com/ivch/dynasty/server/handlers/health"
	svcReqs "github.com/ivch/dynasty/server/handlers/requests"
	repoReqs "github.com/ivch/dynasty/server/handlers/requests/repo"
	transportReqs "github.com/ivch/dynasty/server/handlers/requests/transport"
	transportUI "github.com/ivch/dynasty/server/handlers/ui/transport"
	svcUsers "github.com/ivch/dynasty/server/handlers/users"
	repoUsers "github.com/ivch/dynasty/server/handlers/users/repo"
	transportUsers "github.com/ivch/dynasty/server/handlers/users/transport"
)

// nolint: funlen
func main() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		if err := godotenv.Load(".env"); err != nil {
			stdLog.Fatal("error loading .env file:" + err.Error())
		}
	}

	cfg, err := config.New()
	if err != nil {
		stdLog.Fatal("failed to init config: " + err.Error())
	}

	lvl, err := logger.ParseLevel(cfg.LogLevel)
	if err != nil {
		stdLog.Fatal("failed to create logger: " + err.Error())
	}

	log := logger.NewStdLog(logger.WithLevel(lvl))
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database, cfg.DB.SSL))
	if err != nil {
		stdLog.Fatalf("cannot connect to db: %s", err.Error())
	}

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.S3.Key, cfg.S3.Secret, ""),
		Endpoint:    aws.String(cfg.S3.Endpoint),
		Region:      aws.String(cfg.S3.Region),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		stdLog.Fatalf("cannot start DO session: %s", err)
	}
	s3Client := s3.New(newSession)
	p := bluemonday.StrictPolicy()

	mailSender := email.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Pass, cfg.SMTP.From)

	healthChecker := health.NewMultiChecker()
	healthTransport := health.NewHTTPTransport(healthChecker)

	userService := svcUsers.New(log, repoUsers.New(db), cfg.UserService.VerifyRegCode, cfg.UserService.MembersLimit, mailSender)
	usersTransport := transportUsers.NewHTTPTransport(log, userService, p)
	authService := svcAuth.New(log, repoAuth.New(db), clientUsers.New(userService), cfg.AuthService.JWTSecret)
	authTransport := transportAuth.NewHTTPTransport(log, authService)
	dictService := svcDict.New(log, repoDict.New(db))
	dictTransport := transportDict.NewHTTPTransport(log, dictService)
	reqsSvc := svcReqs.New(log, repoReqs.New(db), s3Client, cfg.RequestService.S3SpaceName, cfg.RequestService.CDNHost)
	reqsTransport := transportReqs.NewHTTPTransport(log, reqsSvc, p)
	uiTransport := transportUI.NewHTTPHandler(cfg.GuardUI.APIHost, cfg.GuardUI.PageURI, cfg.GuardUI.PagerLimit)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-signals
		cancel()
		log.Info("shutdown signal '%s' received! Bye!", sig)
	}()

	srv, err := server.New(
		fmt.Sprintf(":%s", cfg.HTTPPort),
		log,
		map[string]http.Handler{
			"/health":     healthTransport,
			"/users":      usersTransport,
			"/auth":       authTransport,
			"/dictionary": dictTransport,
			"/requests":   reqsTransport,
			"/ui":         uiTransport,
		})
	if err != nil {
		stdLog.Fatal(fmt.Errorf("failed to create server: %w \n", err))
	}

	log.Info("server started to listen on :%s", cfg.HTTPPort)
	if err := srv.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		stdLog.Fatal(fmt.Errorf("server failed: %w", err))
	}
}
