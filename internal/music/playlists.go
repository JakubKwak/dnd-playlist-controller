package music

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Playlist struct {
	Id  int    `json:"id"`
	Uri string `json:"uri"`
}

const playlistsFile = "playlists.json"

// hello i am under the water please help me
func loadPlaylistURIs() (map[int]string, error) {
	var playlists []Playlist
	file, err := os.Open(playlistsFile)
	if err != nil {
		return nil, fmt.Errorf("could not open playlists file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read playlists file: %w", err)
	}

	if err = json.Unmarshal(data, &playlists); err != nil {
		return nil, fmt.Errorf("the json's fucked you dumbass: %w", err)
	}

	uris := make(map[int]string, len(playlists))
	for _, playlist := range playlists {
		uris[playlist.Id] = playlist.Uri
	}
	return uris, nil
}
