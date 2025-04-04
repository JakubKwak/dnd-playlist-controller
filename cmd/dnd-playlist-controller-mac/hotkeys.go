package main

import (
	"fmt"

	"golang.design/x/hotkey"
)

// OS specific key mapping

func getModMap() map[string]hotkey.Modifier {
	return map[string]hotkey.Modifier{
		"OPTION": hotkey.ModOption,
		"CTRL":   hotkey.ModCtrl,
		"SHIFT":  hotkey.ModShift,
		"CMD":    hotkey.ModCmd,
	}
}

func getKeyMap() map[string]hotkey.Key {
	keyMap := map[string]hotkey.Key{}
	// Add number keys (0-9)
	for k, v := range map[string]hotkey.Key{
		"0": hotkey.Key0, "1": hotkey.Key1, "2": hotkey.Key2, "3": hotkey.Key3, "4": hotkey.Key4,
		"5": hotkey.Key5, "6": hotkey.Key6, "7": hotkey.Key7, "8": hotkey.Key8, "9": hotkey.Key9,
	} {
		keyMap[k] = v
	}

	// Add letter keys (A-Z)
	for k, v := range map[string]hotkey.Key{
		"A": hotkey.KeyA, "B": hotkey.KeyB, "C": hotkey.KeyC, "D": hotkey.KeyD, "E": hotkey.KeyE, "F": hotkey.KeyF,
		"G": hotkey.KeyG, "H": hotkey.KeyH, "I": hotkey.KeyI, "J": hotkey.KeyJ, "K": hotkey.KeyK, "L": hotkey.KeyL,
		"M": hotkey.KeyM, "N": hotkey.KeyN, "O": hotkey.KeyO, "P": hotkey.KeyP, "Q": hotkey.KeyQ, "R": hotkey.KeyR,
		"S": hotkey.KeyS, "T": hotkey.KeyT, "U": hotkey.KeyU, "V": hotkey.KeyV, "W": hotkey.KeyW, "X": hotkey.KeyX,
		"Y": hotkey.KeyY, "Z": hotkey.KeyZ,
	} {
		keyMap[k] = v
	}

	// Add function keys (F1-F20)
	for i := 1; i <= 20; i++ {
		keyMap[fmt.Sprintf("F%d", i)] = hotkey.Key(0x7A + (i - 1)) // Uses the provided function key mapping
	}
	return keyMap
}
