package server

import (
	"fmt"
	"github.com/J4stEu/getBlock/internal/app/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// Server - server_errors structure
type Server struct {
	config *config.Config
	logger *logrus.Logger
	router *mux.Router
}

// New - new server_errors instance
func New(config *config.Config, logger *logrus.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
	}
}

// Start - start server_errors instance
func (srv *Server) Start() error {
	srv.ConfigureLogger()
	srv.ConfigureRouter()
	srv.logger.Info("Starting application...")
	return http.ListenAndServe(
		fmt.Sprintf("%s:%s", srv.config.Server.ServerAddr, strconv.Itoa(int(srv.config.Server.ServerPort))), srv.router)
}
