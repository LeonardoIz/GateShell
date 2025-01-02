package proxy

import (
	"fmt"
	"log"
	"net"

	"gateshell/internal/auth"
	"gateshell/internal/utils"

	"golang.org/x/crypto/ssh"
)

type ServerConfig struct {
	Port          string
	HostKeyFile   string
	ServerName    string
	ServerVersion string
	Config        *utils.Manager
}

type Server struct {
	config        *ServerConfig
	sshConfig     *ssh.ServerConfig
	router        *Router
	authenticator *auth.Authenticator
}

func NewServer(config *ServerConfig) *Server {
	return &Server{
		config: config,
		router: NewRouter(config.Config),
		authenticator: auth.NewAuthenticator(&auth.AuthConfig{
			ServerVersion: fmt.Sprintf("SSH-2.0-%s-%s", config.ServerName, config.ServerVersion),
			HostKeyFile:   config.HostKeyFile,
		}),
	}
}

func (s *Server) setupSSHConfig() error {
	// Setup SSH server configuration
	var err error
	s.sshConfig, err = s.authenticator.ConfigureServer()
	return err
}

func (s *Server) handleConnection(nConn net.Conn) {
	// Handle new incoming SSH connection
	conn, chans, reqs, err := ssh.NewServerConn(nConn, s.sshConfig)
	if err != nil {
		log.Printf("Failed to establish SSH connection: %v\n", err)
		return
	}
	defer conn.Close()

	log.Printf("New SSH connection from %s with username %s\n", conn.RemoteAddr(), conn.User())

	password := conn.Permissions.Extensions["password"]
	s.router.HandleConnection(conn, chans, reqs, conn.User(), password)
}

func (s *Server) Start() error {
	// Start the SSH server
	if err := s.setupSSHConfig(); err != nil {
		return fmt.Errorf("failed to setup SSH config: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.config.Port, err)
	}

	log.Printf("Server listening on port %s...\n", s.config.Port)

	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}

		go s.handleConnection(nConn)
	}
}
