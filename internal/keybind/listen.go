package keybind

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func Register(keys []*Keybind) {
	mainthread.Init(func() {
		registerHotkeys(keys)
	})
}

func registerHotkeys(keys []*Keybind) {
	wg := sync.WaitGroup{}
	wg.Add(len(keys))
	for _, key := range keys {
		go func() {
			defer wg.Done()
			err := listenHotkey(key)
			if err != nil {
				fmt.Printf("Error registering hotkey (%s): %s", key.FriendlyName, err)
			}
		}()
	}

	wg.Wait()
}

func listenHotkey(key *Keybind) error {
	hk := hotkey.New(key.Modifiers, key.Key)

	err := hk.Register()
	if err != nil {
		return nil
	}
	defer hk.Unregister()
	fmt.Printf("Listening for hotkey: %s\n", key.FriendlyName)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-hk.Keydown():
			key.Handle()
		case <-sigChan: // If exit signal is received, exit cleanly
			return nil
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
