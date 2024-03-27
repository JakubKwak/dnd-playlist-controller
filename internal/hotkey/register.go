package hotkey

import (
	"fmt"
	"syscall"
)

func Register(user32 *syscall.DLL, hwnd uintptr) (map[int]*Hotkey, error) {
	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Hotkeys to listen to:
	// (hardcoded for now)
	keys := map[int]*Hotkey{
		6:  {6, ModAlt + ModCtrl + ModShift, '6'},
		7:  {7, ModAlt + ModCtrl + ModShift, '7'},
		8:  {8, ModAlt + ModCtrl + ModShift, '8'},
		9:  {9, ModAlt + ModCtrl + ModShift, '9'},
		10: {10, ModAlt + ModCtrl + ModShift, '0'},
	}

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(hwnd, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 != 1 {
			return nil, fmt.Errorf("i shat myself %#v: %w", v, err)
		}
	}
	return keys, nil
}
