package keybind

import (
	"fmt"
	"strings"

	"golang.design/x/hotkey"
)

type Keybind struct {
	Key          hotkey.Key
	Modifiers    []hotkey.Modifier
	Handle       func()
	FriendlyName string
}

// Create a new keybind. keyMap and modMap have to be passed in, as they can differ between operating systems
func NewKeybind(
	key string,
	modifiers []string,
	keyMap map[string]hotkey.Key,
	modMap map[string]hotkey.Modifier,
	handler func(),
) (Keybind, error) {
	friendlyName := ""

	parsedKey, valid := keyMap[strings.ToUpper(key)]
	if !valid {
		return Keybind{}, fmt.Errorf("%s is not a valid key", key)
	}
	mods := []hotkey.Modifier{}
	for _, modifier := range modifiers {
		mod, valid := modMap[strings.ToUpper(modifier)]
		if !valid {
			return Keybind{}, fmt.Errorf("%s is not a valid modifier", modifier)
		}
		friendlyName += " " + strings.ToUpper(modifier)
		mods = append(mods, mod)
	}

	return Keybind{
		Key:          parsedKey,
		Modifiers:    mods,
		Handle:       handler,
		FriendlyName: key + friendlyName,
	}, nil
}
