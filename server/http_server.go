package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"

	"github.com/ivch/dynasty/common/logger"
	"github.com/ivch/dynasty/server/middlewares"
)

// Server represents http server which holds all dependencies.
type Server struct {
	server *http.Server
	log    logger.Logger
	router chi.Router
	health http.Handler
	svc    http.Handler
}

// New returns a new instance of Server struct.
func New(addr string, log logger.Logger, services map[string]http.Handler) (*Server, error) {
	router := chi.NewRouter()
	server := &http.Server{
		Handler: router,
		Addr:    addr,
	}

	s := &Server{
		router: router,
		server: server,
		log:    log,
	}
	return s.routes(services)
}

func (s *Server) routes(services map[string]http.Handler) (*Server, error) {
	logmw := middlewares.NewLogging(s.log)
	recmw := middlewares.NewRecover(s.log)
	idctxmw := middlewares.NewIDCtx(s.log)
	s.router.Use(recmw.Middleware, idctxmw.Middleware, chimw.StripSlashes, chimw.RequestID)
	for prefix, service := range services {
		s.router.With(logmw.Middleware).Mount(prefix, service)
	}

	return s, nil
}

func (s *Server) Serve(ctx context.Context) error {
	// handle shutdown signal in background
	go s.handleShutdown(ctx)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) handleShutdown(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if err := s.shutdown(); err != nil {
				s.log.Error("killing server!", err)
				os.Exit(1)
			}
		default:
			continue
		}
	}
}

func (s *Server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown the server: %w", err)
	}
	return nil
}
