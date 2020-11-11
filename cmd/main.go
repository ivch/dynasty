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

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	// uCli "github.com/ivch/dynasty/common/clients/users"
	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/config"
	"github.com/ivch/dynasty/server"
	"github.com/ivch/dynasty/server/handlers/health"
	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/handlers/users/repo"
	"github.com/ivch/dynasty/server/handlers/users/transport"
)

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
		log.Error("cannot connect to db: %s", err.Error())
	}

	// s3Config := &aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(cfg.S3.Key, cfg.S3.Secret, ""),
	// 	Endpoint:    aws.String(cfg.S3.Endpoint),
	// 	Region:      aws.String(cfg.S3.Region),
	// }

	// newSession, err := session.NewSession(s3Config)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("cannot start do session")
	// }
	// s3Client := s3.New(newSession)
	p := bluemonday.StrictPolicy()

	healthChecker := health.NewMultiChecker()
	healthTransport := health.NewHTTPTransport(healthChecker)

	userService := users.New(log, repo.NewUsers(db), cfg.UserService.VerifyRegCode, cfg.UserService.MembersLimit)
	usersTransport := transport.NewHTTPTransport(log, userService, p)
	// authModule, _ := auth.New(log, repository.NewAuth(db), uCli.New(userService), cfg.AuthService.JWTSecret)
	// requestsModule, _ := requests.New(logger, repository.NewRequests(db), p, s3Client, cfg.RequestService.S3SpaceName, cfg.RequestService.CDNHost)
	// dictionariesModule, _ := dictionaries.New(repository.NewDictionaries(db), logger)
	//
	// r := chi.NewRouter()
	// r.Use(accessLogMiddleware(logger))
	//
	// r.Mount("/users", usersModule)
	// r.Mount("/auth", authModule)
	// r.Mount("/requests", requestsModule)
	// r.Mount("/dictionary", dictionariesModule)
	// r.Mount("/ui", ui.NewHTTPHandler(cfg.GuardUI.APIHost, cfg.GuardUI.PageURI, cfg.GuardUI.PagerLimit))
	//
	// r.Get("/health", func(w http.ResponseWriter, r *http.Request) {})
	//

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
			"/health": healthTransport,
			"/users":  usersTransport,
		})
	if err != nil {
		stdLog.Fatal(fmt.Errorf("failed to create server: %w \n", err))
	}

	log.Info("server started to listen on :%s", cfg.HTTPPort)
	if err := srv.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		stdLog.Fatal(fmt.Errorf("server failed: %w", err))
	}
}

// func newLogger(verbose bool) (logger *zerolog.Logger) {
// 	switch verbose {
// 	case true:
// 		devLogger := zerolog.New(zerolog.ConsoleWriter{
// 			NoColor: false,
// 			Out:     os.Stdout,
// 		}).Level(zerolog.DebugLevel).With().Timestamp().Logger()
// 		logger = &devLogger
// 	default:
// 		prodLogger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()
// 		logger = &prodLogger
// 	}
// 	return logger
// }

// var accessLogMiddleware = func(log *zerolog.Logger) func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()
// 			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
// 			next.ServeHTTP(ww, r)
// 			duration := time.Since(start)
// 			log.Info().
// 				Str("tag", "http_log").
// 				Str("remote_addr", r.RemoteAddr).
// 				Str("request_method", r.Method).
// 				Str("request_uri", r.RequestURI).
// 				Int("response_code", ww.Status()).
// 				Dur("duration", duration).
// 				Msg("request")
// 		})
// 	}
// }
