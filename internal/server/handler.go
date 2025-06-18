package server

import "github.com/xeaser/citriage/internal/buildsapi"

func (s *Server) registerHandlers() {
	buildsapi.RegisterHandlers(s.router, s.cfg, s.logger)
}
