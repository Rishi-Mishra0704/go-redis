package main

import "net"

const defaultListenAddress = ":5001"

// Config is the configuration for the server
type Config struct {
	// The port to listen on
	listenAddress string
}

// Server represents a server that listens for incoming connections
type Server struct {
	Config
	ln net.Listener
}

// NewServer creates a new server with the given configuration
func NewServer(config Config) *Server {
	if len(config.listenAddress) == 0 {
		config.listenAddress = defaultListenAddress
	}
	return &Server{
		Config: config,
	}
}

// Start starts the server and listens for incoming connections using the tcp protocol
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	return nil
}
