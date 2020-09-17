package spotify

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/flexicon/spotimoods-go/internal"
)

// do performs an action against the spotify API and on authorization failure attempts to refresh the given token and try again
func (c *Client) do(req *http.Request, token *internal.SpotifyToken) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+token.Token)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	logBodyAndRewind(resp)

	if resp.StatusCode != http.StatusUnauthorized || token.Refresh == "" {
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			return nil, httpStatusErr(resp)
		}
		return resp, nil
	}
	resp.Body.Close()

	if err := c.Refresh(token); err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	resp, err = c.http.Do(req)
	if err != nil {
		return nil, err
	}
	logBodyAndRewind(resp)

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, httpStatusErr(resp)
	}

	return resp, nil
}

// fetch the given request and return the raw response body
func (c *Client) fetch(req *http.Request, token *internal.SpotifyToken) ([]byte, error) {
	resp, err := c.do(req, token)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// fetch the given request from cache, execute it otherwise and return the raw response body
func (c *Client) fetchWithCache(req *http.Request, token *internal.SpotifyToken, cacheItem *internal.CacheItem) ([]byte, error) {
	if c.cache.Exists(cacheItem.Key) {
		var body []byte
		if err := c.cache.Get(cacheItem.Key, &body); err != nil {
			return nil, err
		}
		return body, nil
	}

	// Execute actual request
	body, err := c.fetch(req, token)
	if err != nil {
		return nil, err
	}

	// Cache results
	cacheItem.Value = &body
	c.cache.Set(cacheItem)

	return body, nil
}

func httpStatusErr(resp *http.Response) error {
	body, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("Http status %d: %s", resp.StatusCode, body)
}

func logBodyAndRewind(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	resp.Body.Close()
	log.Printf("HTTP Response body: %s", body)

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return nil
}
