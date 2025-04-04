package main

import (
	"dnd-playlist-controller/internal/keybind"
	"dnd-playlist-controller/internal/playlistswitcher"
	"dnd-playlist-controller/internal/spotifyclient"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(".env"); err != nil {
		error(fmt.Sprintf("Could not load env: %s", err))
	}

	// create spotify client
	client, err := spotifyclient.New()
	if err != nil {
		error(err.Error())
	}

	// create the playlist switcher
	switcher := playlistswitcher.NewSwitcher(client)

	// load playlist URIs and their hotkeys from JSON
	hotkeys, err := playlistswitcher.InitPlaylistHotkeys(switcher)
	if err != nil {
		error(err.Error())
	}

	keybind.Register(hotkeys)
}

func error(err string) {
	log.Fatalf("(╯°□°）╯︵ ┻━┻ %s", err)
}
