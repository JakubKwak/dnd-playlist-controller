package playlistswitcher

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

type Switcher struct {
	client *spotify.Client
	ctx    context.Context
}

// empty switcher, used for debugging
func FakeSwitcher() *Switcher {
	return &Switcher{ctx: context.Background()}
}

func NewSwitcher(client *spotify.Client) *Switcher {
	return &Switcher{client: client, ctx: context.Background()}
}

func (s *Switcher) SwitchPlaylist(uri string) {
	spotifyUri := (spotify.URI)("spotify:playlist:" + uri)
	fmt.Printf("Switching to playlist: %s\n", spotifyUri)

	err := s.client.PlayOpt(s.ctx, &spotify.PlayOptions{PlaybackContext: &spotifyUri})
	if err != nil {
		fmt.Printf("i shid and fard: %s", err)
	}
}
