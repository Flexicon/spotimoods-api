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
	Mood() MoodRepository
}

// ServiceProvider manages all services
type ServiceProvider struct {
	repos   RepositoryProvider
	spotify SpotifyClient
	queue   QueueService
	cache   Cache
}

// NewServiceProvider constructor
func NewServiceProvider(repos RepositoryProvider, spotify SpotifyClient, qs QueueService, c Cache) *ServiceProvider {
	return &ServiceProvider{
		repos:   repos,
		spotify: spotify,
		queue:   qs,
		cache:   c,
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

// Mood returns a new Mood service
func (p *ServiceProvider) Mood() *MoodService {
	return NewMoodService(p.repos.Mood(), p.Queue(), p.Spotify())
}

// Queue returns the Queue service instance
func (p *ServiceProvider) Queue() QueueService {
	return p.queue
}

// Cache returns the Cache service instance
func (p *ServiceProvider) Cache() Cache {
	return p.cache
}
