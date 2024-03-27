package main

import (
	"dnd-playlist-controller/internal/hotkey"
	"dnd-playlist-controller/internal/music"
	"fmt"
	"log"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(".env"); err != nil {
		error(fmt.Sprintf("Could not load env: %s", err))
	}

	daddy, err := music.NewSwitcher(music.SpotifyClient())
	if err != nil {
		error(err.Error())
	}

	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	keys, err := hotkey.Register(user32)
	if err != nil {
		error(fmt.Sprintf("hot keys in my area: %s", err))
	}
	hotkey.Listen(user32, keys, daddy)
}

func error(err string) {
	log.Fatalf("(╯°□°）╯︵ ┻━┻ %s", err)
}
