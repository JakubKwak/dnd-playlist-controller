package main

import (
	"fmt"

	"golang.design/x/hotkey"
)

// OS specific key mapping

func getModMap() map[string]hotkey.Modifier {
	return map[string]hotkey.Modifier{
		"ALT":   hotkey.ModAlt,
		"CTRL":  hotkey.ModCtrl,
		"SHIFT": hotkey.ModShift,
		"WIN":   hotkey.ModWin,
	}
}

func getKeyMap() map[string]hotkey.Key {
	keyMap := map[string]hotkey.Key{}
	// Add numbers 0-9
	for i := 0; i <= 9; i++ {
		keyMap[fmt.Sprintf("%d", i)] = hotkey.Key(0x30 + i)
	}
	// Add letters A-Z
	for i := 'A'; i <= 'Z'; i++ {
		keyMap[string(i)] = hotkey.Key(i)
	}
	// Add function keys F1-F20
	for i := 1; i <= 20; i++ {
		keyMap[fmt.Sprintf("F%d", i)] = hotkey.Key(0x70 + (i - 1))
	}
	return keyMap
}
