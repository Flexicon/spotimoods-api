package spotify

// PlaylistPayload for
// https://developer.spotify.com/documentation/web-api/reference/playlists/create-playlist/
type PlaylistPayload struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
