package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/flexicon/spotimoods-go/internal"
	"github.com/pkg/errors"
)

// GetArtistsByIDs retrieves the artists related to the given IDs
func (c *Client) GetArtistsByIDs(token *internal.SpotifyToken, ids []string) ([]*internal.SpotifyArtist, error) {
	// First check cache for artists and only make a request if any non-cached artists remain
	artists, remainingIDs := c.getAndFilterCachedArtists(ids)
	if len(remainingIDs) == 0 {
		return artists, nil
	}

	artistsURL, err := url.Parse(fmt.Sprintf("https://api.spotify.com/v1/artists?ids=%s", strings.Join(remainingIDs, ",")))
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare artists by id url")
	}

	req, _ := http.NewRequest(http.MethodGet, artistsURL.String(), nil)
	body, err := c.fetch(req, token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve artists by ids")
	}

	var response struct {
		Artists []*internal.SpotifyArtist `json:"artists"`
	}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing artists response: %v", err)
	}
	c.cacheArtists(response.Artists)

	return append(artists, response.Artists...), nil
}

// getAndFilterCachedArtists tries to retrieve each artist from cache by id,
// returns a slice of artists and a slice of remaining filtered ids that weren't found in cache
func (c *Client) getAndFilterCachedArtists(ids []string) ([]*internal.SpotifyArtist, []string) {
	artists := make([]*internal.SpotifyArtist, 0)
	remaining := make([]string, 0)

	for _, id := range ids {
		var artist *internal.SpotifyArtist
		err := c.cache.Get(fmt.Sprintf("Artist-%s", id), &artist)

		if err == nil && artist != nil {
			artists = append(artists, artist)
		} else {
			remaining = append(remaining, id)
		}
	}

	return artists, remaining
}

func (c *Client) cacheArtists(artists []*internal.SpotifyArtist) {
	for _, artist := range artists {
		err := c.cache.Set(&internal.CacheItem{
			Key:   fmt.Sprintf("Artist-%s", artist.ID),
			Value: &artist,
			TTL:   time.Minute * 15,
		})

		if err != nil {
			log.Println(errors.Wrap(err, "failed to store artist in cache"))
		}
	}
}
