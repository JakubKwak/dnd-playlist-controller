package hotkey

import (
	"fmt"
	"syscall"
)

func Register(user32 *syscall.DLL, hwnd uintptr, keys map[int]*Hotkey) error {
	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(hwnd, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 != 1 {
			return fmt.Errorf("error registering hotkey: %#v: %w", v, err)
		}
	}
	return nil
}
