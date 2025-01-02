// internal/proxy/router.go
package proxy

import (
	"fmt"
	"io"
	"log"

	"gateshell/internal/auth"
	"gateshell/internal/config"

	"golang.org/x/crypto/ssh"
)

type Router struct {
	authenticator *auth.Authenticator
	config        *config.Manager
}

func NewRouter(config *config.Manager) *Router {
	return &Router{
		authenticator: auth.NewAuthenticator(&auth.AuthConfig{}),
		config:        config,
	}
}

func (r *Router) connectToUpstream(username, password string) (*ssh.Client, error) {
	endpoint, err := r.config.GetEndpoint(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint configuration: %v", err)
	}

	sshConfig := r.authenticator.GetUpstreamConfig(endpoint.Auth.User, password)

	client, err := ssh.Dial("tcp", endpoint.Target, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial upstream server: %v", err)
	}

	return client, nil
}

func (r *Router) handleChannels(chans <-chan ssh.NewChannel, upstreamClient *ssh.Client) {
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		downstreamChannel, downstreamRequests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Failed to accept downstream channel: %v", err)
			continue
		}

		upstreamChannel, upstreamRequests, err := upstreamClient.OpenChannel("session", nil)
		if err != nil {
			log.Printf("Failed to open upstream channel: %v", err)
			downstreamChannel.Close()
			continue
		}

		go ssh.DiscardRequests(upstreamRequests)
		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok, err := upstreamChannel.SendRequest(req.Type, req.WantReply, req.Payload)
				if err != nil {
					log.Printf("Failed to forward request: %v", err)
					continue
				}
				if req.WantReply {
					req.Reply(ok, nil)
				}
			}
		}(downstreamRequests)

		go func() {
			_, _ = io.Copy(upstreamChannel, downstreamChannel)
			upstreamChannel.Close()
		}()
		go func() {
			_, _ = io.Copy(downstreamChannel, upstreamChannel)
			downstreamChannel.Close()
		}()
	}
}

func (r *Router) HandleConnection(conn *ssh.ServerConn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request, username, password string) {
	upstreamClient, err := r.connectToUpstream(username, password)
	if err != nil {
		log.Printf("Failed to connect to upstream server: %v\n", err)
		return
	}
	defer upstreamClient.Close()

	go ssh.DiscardRequests(reqs)
	r.handleChannels(chans, upstreamClient)
}
