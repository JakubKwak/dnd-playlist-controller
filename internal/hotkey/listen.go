package hotkey

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

type playlistSwitcher interface {
	Switch(int)
}

func Listen(user32 *syscall.DLL, keys map[int]*Hotkey, switcher playlistSwitcher) {
	fmt.Println("( ͡° ͜ʖ ͡°) Listening for hotkeys...")

	peekmsg := user32.MustFindProc("PeekMessageW")

	for {
		var msg = &MSG{}
		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Yooooooo a hotkey was pressed:", keys[int(id)])
			if id == 6 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				return
			}

			switcher.Switch(int(id))
		}
		fmt.Printf("%#v \n", msg)
		time.Sleep(time.Millisecond * 50)
	}
}
