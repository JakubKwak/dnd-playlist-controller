package main

import (
	"dnd-playlist-controller/internal/hotkey"
	"dnd-playlist-controller/internal/playlistswitcher"
	"dnd-playlist-controller/internal/spotifyclient"
	"fmt"
	"log"
	"runtime"
	"syscall"

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

	// load playlist URIs and their hotkeys from JSON
	playlistURIs, hotkeys, err := playlistswitcher.LoadPlaylistHotkeys()
	if err != nil {
		error(err.Error())
	}

	// create the playlist switcher
	switcher := playlistswitcher.NewSwitcher(client, playlistURIs)
	// switcher := music.FakeSwitcher()

	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	// lock thread, as registering hotkeys is thread-specific
	runtime.LockOSThread()
	hwnd, err := hotkey.GiveSimpleWindowPls(user32)
	if err != nil {
		error(err.Error())
	}
	if err = hotkey.Register(user32, hwnd, hotkeys); err != nil {
		error(fmt.Sprintf("hotkey cringe %s", err))
	}
	hotkey.Listen(user32, hotkeys, switcher, hwnd)
}

func error(err string) {
	log.Fatalf("(╯°□°）╯︵ ┻━┻ %s", err)
}
