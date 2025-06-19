package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/cache"
)

const (
	serverDefaultPort = 8080
)

type Server struct {
	cfg        *config.Config
	logger     *slog.Logger
	router     *mux.Router
	handler    http.Handler
	downloader *cache.Downloader
}

func New(cfg *config.Config, logger *slog.Logger, downloader *cache.Downloader) *Server {
	return &Server{
		cfg:        cfg,
		logger:     logger,
		router:     mux.NewRouter(),
		downloader: downloader,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) {
	s.registerHandlers()
	s.handler = s.router

	if err := s.printRoutes(); err != nil {
		s.logger.WarnContext(ctx, fmt.Sprintf("Unable to print server routes, err: %v", err))
	}

	s.StartServer(ctx)
}

func (s *Server) StartServer(ctx context.Context) *http.Server {
	if s.cfg.Server.Port == 0 {
		s.logger.ErrorContext(ctx, "Server port not in config, using default port")
		s.cfg.Server.Port = serverDefaultPort
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler: s.handler,
	}

	go func() {
		s.logger.InfoContext(ctx, fmt.Sprintf("Starting HTTP Server on Port: %d", s.cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Unable to start HTTP server, err : %v", err)
		}
	}()

	return srv
}

func (s *Server) printRoutes() error {
	s.logger.Info("Routes Configured on HTTP Server")
	return s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		s.logger.Info(fmt.Sprintf("Path: %v, Method: %v", pathTemplate, methods))
		return nil
	})
}
