package config

import (
	"github.com/J4stEu/getBlock/internal/app/errors"
	"github.com/J4stEu/getBlock/internal/app/errors/server_errors"
	"github.com/J4stEu/getBlock/internal/pkg"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// Server - server_errors configuration
type Server struct {
	ServerAddr string
	ServerPort uint
	LogLevel   string
	APIkey     string
}

// Config - application configuration
type Config struct {
	Server *Server
}

func CheckENV() bool {
	_, err := os.LookupEnv("SERVER_ADDR")
	if !err {
		return false
	}
	_, err = os.LookupEnv("SERVER_PORT")
	if !err {
		return false
	}
	_, err = os.LookupEnv("LOG_LEVEL")
	if !err {
		return false
	}
	_, err = os.LookupEnv("API_KEY")
	if !err {
		return false
	}
	return true
}

func ReadConfiguration(logger *logrus.Logger) *Config {
	config := &Config{&Server{}}

	// Server configuration
	// ServerAddr
	serverAddr, err := os.LookupEnv("SERVER_ADDR")
	if !err {
		logger.Fatal(errors.ServerErrorLevel, server_errors.EnvReadError)
	}
	validServerAddr := pkg.IsValidIP(serverAddr)
	if !validServerAddr {
		logger.WithFields(log.Fields{
			"error": "Invalid server IP address.",
		}).Fatal(errors.SetError(errors.ServerErrorLevel, server_errors.EnvSetError))
	}
	// ServerPort
	config.Server.ServerAddr = serverAddr
	serverPort, err := os.LookupEnv("SERVER_PORT")
	if !err {
		logger.Fatal(errors.ServerErrorLevel, server_errors.EnvReadError)
	}
	serverPortUINT, convertErr := strconv.Atoi(serverPort)
	if convertErr != nil {
		logger.WithFields(log.Fields{
			"error": convertErr,
		}).Fatal(errors.SetError(errors.ServerErrorLevel, server_errors.EnvSetError))
	}
	config.Server.ServerPort = uint(serverPortUINT)
	// LogLevel
	logLevel, err := os.LookupEnv("LOG_LEVEL")
	if !err {
		logger.Fatal(errors.ServerErrorLevel, server_errors.EnvReadError)
	}
	config.Server.LogLevel = logLevel
	// APIkey
	apiKey, err := os.LookupEnv("API_KEY")
	if !err {
		logger.Fatal(errors.ServerErrorLevel, server_errors.EnvReadError)
	}
	config.Server.APIkey = apiKey

	return config
}

func DefaultConfiguration() *Config {
	return &Config{
		Server: &Server{
			ServerAddr: "localhost",
			ServerPort: 8080,
			LogLevel:   "debug",
			APIkey:     "2687bcf0-cd62-4de6-9d4e-13606e872040",
		},
	}
}
