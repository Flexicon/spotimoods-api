package internal

import (
	"net/http"
)

// HTTPClient for making outbound requests
type HTTPClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// RepositoryProvider manages all repositories
type RepositoryProvider interface {
	User() UserRepository
}

// ServiceProvider manages all services
type ServiceProvider struct {
	repos   RepositoryProvider
	spotify SpotifyClient
}

// NewServiceProvider constructor
func NewServiceProvider(repos RepositoryProvider, spotify SpotifyClient) *ServiceProvider {
	return &ServiceProvider{
		repos:   repos,
		spotify: spotify,
	}
}

// User returns a new User service
func (p *ServiceProvider) User() *UserService {
	return NewUserService(p.repos.User())
}

// Spotify returns a new Spotify client
func (p *ServiceProvider) Spotify() SpotifyClient {
	return p.spotify
}
