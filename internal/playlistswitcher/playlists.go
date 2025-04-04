package playlistswitcher

import (
	"dnd-playlist-controller/internal/hotkey"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type PlaylistHotkey struct {
	Id        int      `json:"id"`
	Uri       string   `json:"uri"`
	Modifiers []string `json:"modifiers"`
	Key       string   `json:"key"`
}

const hotkeysFile = "hotkeys.json"

func LoadPlaylistHotkeys() (map[int]string, map[int]*hotkey.Hotkey, error) {
	var playlists []PlaylistHotkey
	file, err := os.Open(hotkeysFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open playlists file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read playlists file: %w", err)
	}

	if err = json.Unmarshal(data, &playlists); err != nil {
		return nil, nil, fmt.Errorf("the json's fucked you dumbass: %w", err)
	}

	keys := make(map[int]*hotkey.Hotkey, len(playlists))
	uris := make(map[int]string, len(playlists))
	for _, playlist := range playlists {
		mods, err := addModifiers(playlist.Modifiers)
		if err != nil {
			return nil, nil, err
		}
		key, err := parseKey(playlist.Key)
		if err != nil {
			return nil, nil, err
		}

		uris[playlist.Id] = playlist.Uri
		keys[playlist.Id] = &hotkey.Hotkey{
			Id:        playlist.Id,
			Modifiers: mods,
			KeyCode:   int(key),
		}
	}
	return uris, keys, nil
}

func parseKey(key string) (rune, error) {
	runeArray := []rune(key)
	if len(runeArray) != 1 {
		return 0, fmt.Errorf("the key in hotkeys.json must be a single character")
	}
	return runeArray[0], nil
}

func addModifiers(mods []string) (int, error) {
	var total int
	for _, mod := range mods {
		val, ok := hotkey.ParseMod(mod)
		if !ok {
			return 0, fmt.Errorf("%s is not a valid modifier", mod)
		}
		total += val
	}
	return total, nil
}
