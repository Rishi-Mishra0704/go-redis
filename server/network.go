package server

import (
	"log/slog"
	"net"
)

const defaultListenAddress = ":5001"

// Config is the configuration for the server
type Config struct {
	// The port to listen on
	listenAddress string
}

// Server represents a server that listens for incoming connections
type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
}

// NewServer creates a new server with the given configuration
func NewServer(config Config) *Server {
	if len(config.listenAddress) == 0 {
		config.listenAddress = defaultListenAddress
	}
	return &Server{
		Config:    config,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
	}
}

// Start starts the server and listens for incoming connections using the tcp protocol
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()

	return s.acceptLoop()

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept failed: %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer
	peer.readLoop()
}
