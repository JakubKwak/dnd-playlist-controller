package main

import (
	"dnd-playlist-controller/internal/keybind"
	"dnd-playlist-controller/internal/playlistswitcher"
	"dnd-playlist-controller/internal/spotifyclient"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	chdirToExecutableIfNeeded()

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

// chdirToExecutableIfNeeded switches the working directory to the executable's directory
// when hotkeys.json isn't found in the current working directory (e.g. when double-clicking the binary).
func chdirToExecutableIfNeeded() {
	if _, err := os.Stat("hotkeys.json"); err == nil {
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	exeDir := filepath.Dir(exePath)
	if _, err := os.Stat(filepath.Join(exeDir, "hotkeys.json")); err == nil {
		_ = os.Chdir(exeDir)
	}
}

func error(err string) {
	log.Fatalf("(╯°□°）╯︵ ┻━┻ %s", err)
}
