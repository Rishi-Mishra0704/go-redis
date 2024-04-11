package server

import (
	"fmt"
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
	peers      map[*Peer]bool
	ln         net.Listener
	addPeerCh  chan *Peer
	quitPeerCh chan struct{}
	msgCh      chan []byte
}

// NewServer creates a new server with the given configuration
func NewServer(config Config) *Server {
	if len(config.listenAddress) == 0 {
		config.listenAddress = defaultListenAddress
	}
	return &Server{
		Config:     config,
		peers:      make(map[*Peer]bool),
		addPeerCh:  make(chan *Peer),
		quitPeerCh: make(chan struct{}),
		msgCh:      make(chan []byte),
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
	slog.Info("server is listening", "listenAddr", s.listenAddress)

	return s.acceptLoop()

}

func (s *Server) HandleRawMessage(rawMsg []byte) error {
	fmt.Println(string(rawMsg))
	return nil
}
func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh:
			if err := s.HandleRawMessage(rawMsg); err != nil {
				slog.Error("server is listening", "listenAddr", s.listenAddress)

			}
			fmt.Println((rawMsg))
		case <-s.quitPeerCh:
			return
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
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("New Peer Connected", "RemoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("Error Reading: %v", err)
	}
}
