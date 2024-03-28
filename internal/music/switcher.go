package music

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

type Switcher struct {
	playlistURIs map[int]string
	client       *spotify.Client
	ctx          context.Context
}

// create empty switcher for debug shit
func PoopSwither() *Switcher {
	uris := make(map[int]string, 0)
	return &Switcher{playlistURIs: uris, ctx: context.Background()}
}

func NewSwitcher(client *spotify.Client) (*Switcher, error) {
	playlistURIs, err := loadPlaylistURIs()
	if err != nil {
		return nil, err
	}

	return &Switcher{playlistURIs: playlistURIs, client: client, ctx: context.Background()}, nil
}

func (s *Switcher) Switch(id int) {
	uri, ok := s.playlistURIs[id]
	if !ok {
		return
	}
	spotifyUri := (spotify.URI)("spotify:playlist:" + uri)
	fmt.Printf("Switching to playlist: %s\n", spotifyUri)

	err := s.client.PlayOpt(s.ctx, &spotify.PlayOptions{PlaybackContext: &spotifyUri})
	if err != nil {
		fmt.Printf("i shid and fard: %s", err)
	}
}
