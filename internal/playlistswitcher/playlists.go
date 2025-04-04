package playlistswitcher

import (
	"dnd-playlist-controller/internal/keybind"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.design/x/hotkey"
)

type PlaylistHotkey struct {
	Uri       string   `json:"uri"`
	Modifiers []string `json:"modifiers"`
	Key       string   `json:"key"`
}

const hotkeysFile = "hotkeys.json"

func InitPlaylistHotkeys(switcher *Switcher, keyMap map[string]hotkey.Key, modMap map[string]hotkey.Modifier) ([]*keybind.Keybind, error) {
	var playlists []PlaylistHotkey
	file, err := os.Open(hotkeysFile)
	if err != nil {
		return nil, fmt.Errorf("could not open playlists file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read playlists file: %w", err)
	}

	if err = json.Unmarshal(data, &playlists); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	keys := make([]*keybind.Keybind, 0, len(playlists))
	for _, playlist := range playlists {
		handleFunc := func() {
			switcher.SwitchPlaylist(playlist.Uri)
		}
		hkey, err := keybind.NewKeybind(
			playlist.Key,
			playlist.Modifiers,
			keyMap,
			modMap,
			handleFunc,
		)
		if err != nil {
			return nil, err
		}

		keys = append(keys, &hkey)
	}
	return keys, nil
}
