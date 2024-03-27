package hotkey

import (
	"fmt"
	"syscall"
)

func Register(user32 *syscall.DLL) (map[int]*Hotkey, error) {
	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Hotkeys to listen to:
	// (hardcoded for now)
	keys := map[int]*Hotkey{
		1: {1, ModAlt + ModCtrl + ModShift, '6'},
		2: {2, ModAlt + ModCtrl + ModShift, '7'},
		3: {3, ModAlt + ModCtrl + ModShift, '8'},
		4: {4, ModAlt + ModCtrl + ModShift, '9'},
		5: {5, ModAlt + ModCtrl + ModShift, '0'},
		6: {6, ModAlt + ModCtrl, 'X'}, // ALT+CTRL+X
	}

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 != 1 {
			return nil, fmt.Errorf("i shat myself %#v: %w", v, err)
		}
	}
	return keys, nil
}
