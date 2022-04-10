package server

import (
	"fmt"
	"github.com/J4stEu/getBlock/internal/app/server/handlers"
	"io"
	"net/http"
	"strconv"
)

// ConfigureRouter - router configuration
func (srv *Server) ConfigureRouter() {
	srv.router.HandleFunc("/", srv.HandleRoot())
	srv.router.HandleFunc("/golang_developer_task", srv.HandleTask())
}

// HandleRoot - root route
func (srv *Server) HandleRoot() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		srv.logger.Debug(fmt.Sprintf("%s: %s:%s%s", request.Method, srv.config.Server.ServerAddr, strconv.Itoa(int(srv.config.Server.ServerPort)), request.RequestURI))
		_, err := io.WriteString(writer, "Root")
		if err != nil {
			srv.logger.Warn(err)
		}
	}
}

// HandleTask - golang developer task route
func (srv *Server) HandleTask() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handleTask, err := handlers.HandleTask(srv.config.Server.APIkey)
		if err != nil {
			_, err = io.WriteString(writer, err.Error())
			if err != nil {
				srv.logger.Warn(err)
			}
		}
		_, err = io.WriteString(writer, fmt.Sprintf("Адрес, баланс которого изменился больше остальных (по абсолютной величине) за последние 100 блоков: %s", handleTask))
		if err != nil {
			srv.logger.Warn(err)
		}
	}
}
