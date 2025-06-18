package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xeaser/citriage/config"
)

const (
	serverDefaultPort = 8080
)

type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	router *mux.Router
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		router: mux.NewRouter(),
	}
}

func (s *Server) ListenAndServe(ctx context.Context) {
	s.registerHandlers()

	if err := s.printRoutes(); err != nil {
		s.logger.WarnContext(ctx, "Unable to print server routes", err)
	}

	s.StartServer(ctx)
}

func (s *Server) StartServer(ctx context.Context) *http.Server {
	if s.cfg.ServerConfig.Port == 0 {
		s.logger.ErrorContext(ctx, "Server port not in config, using default port")
		s.cfg.ServerConfig.Port = 8080
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("%v:%v", s.cfg.ServerConfig.Host, s.cfg.ServerConfig.Port),
	}

	go func() {
		s.logger.InfoContext(ctx, fmt.Sprintf("Starting HTTP Server on Port: %d", s.cfg.ServerConfig.Port))
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
		queriesTemplate, err := route.GetQueriesTemplates()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		s.logger.Info(fmt.Sprintf("Path: %v, Method: %v, QueriesTemplate: %v", pathTemplate, methods, queriesTemplate))
		return nil
	})
}
