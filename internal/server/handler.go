package server

import "github.com/xeaser/citriage/internal/logsapi"

func (s *Server) registerHandlers() {
	logsapi.RegisterHandlers(s.router, s.cfg, s.logger, s.downloader)
}
