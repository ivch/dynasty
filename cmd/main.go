package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	uCli "github.com/ivch/dynasty/clients/users"
	"github.com/ivch/dynasty/config"
	"github.com/ivch/dynasty/modules/auth"
	"github.com/ivch/dynasty/modules/users"
)

var (
	// Version is the current version of application
	Version = "0"
	// Branch is the branch this binary was built from
	Branch = "0"
	// Commit is the commit this binary was built from
	Commit = "0"
	// BuildTime is the time this binary was built
	BuildTime = time.Now().Format(time.RFC822)
)

func main() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("error loading .env file:" + err.Error())
		}
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed to init config: " + err.Error())
	}

	logger := newLogger(cfg.LogVerbose)
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database, cfg.DB.SSL))
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to db")
	}

	usersModule, userService := users.New(db, cfg.UserService.VerifyRegCode, logger)
	authModule := auth.New(logger, db, uCli.New(userService), cfg.AuthService.JWTSecret)

	r := chi.NewRouter()
	r.Use(accessLogMiddleware(logger))

	r.Mount("/users", usersModule)
	r.Mount("/auth", authModule)

	r.Get("/about", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"version": Version,
			"branch":  Branch,
			"commit":  Commit,
			"time":    BuildTime,
		}) // nolint: errcheck
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {})

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Fatal().Err(err).Msg("error on server shutdown")
		}

		close(signals)
	}()

	logger.Info().Msg(fmt.Sprintf("HTTP listener started on :%s @ %s", cfg.HTTPPort, time.Now().Format(time.RFC3339)))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err)
	}
}

func newLogger(verbose bool) (logger *zerolog.Logger) {
	switch verbose {
	case true:
		devLogger := zerolog.New(zerolog.ConsoleWriter{
			NoColor: false,
			Out:     os.Stdout,
		}).Level(zerolog.DebugLevel).With().Timestamp().Logger()
		logger = &devLogger
	default:
		prodLogger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()
		logger = &prodLogger
	}
	return logger
}

var accessLogMiddleware = func(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			duration := time.Since(start)
			log.Info().
				Str("tag", "http_log").
				Str("remote_addr", r.RemoteAddr).
				Str("request_method", r.Method).
				Str("request_uri", r.RequestURI).
				Int("response_code", ww.Status()).
				Dur("duration", duration).
				Msg("request")
		})
	}
}
