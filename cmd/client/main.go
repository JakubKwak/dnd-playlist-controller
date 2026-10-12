package main

import (
	"dnd-playlist-controller/internal/keybind"
	"dnd-playlist-controller/internal/playlistswitcher"
	"dnd-playlist-controller/internal/spotifyclient"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Optional: load .env if present (e.g. to override AUTH_SERVER_URL when self-hosting)
	_ = godotenv.Load(".env")

	client, err := spotifyclient.New()
	if err != nil {
		error(err.Error())
	}

	switcher := playlistswitcher.NewSwitcher(client)
	keyMap := getKeyMap()
	modMap := getModMap()

	// load playlist URIs and their hotkeys from JSON
	hotkeys, err := playlistswitcher.InitPlaylistHotkeys(switcher, keyMap, modMap)
	if err != nil {
		error(err.Error())
	}

	keybind.Register(hotkeys)
}

func error(err string) {
	log.Fatalf("(╯°□°）╯︵ ┻━┻ %s", err)
}
