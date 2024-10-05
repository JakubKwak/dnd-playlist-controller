package playlistswitcher

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
func FakeSwitcher() *Switcher {
	uris := make(map[int]string, 0)
	return &Switcher{playlistURIs: uris, ctx: context.Background()}
}

func NewSwitcher(client *spotify.Client, playlistURIs map[int]string) *Switcher {
	return &Switcher{playlistURIs: playlistURIs, client: client, ctx: context.Background()}
}

func (s *Switcher) HandleHotkey(id int) {
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
