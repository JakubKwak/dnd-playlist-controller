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

const WM_HOTKEY = 0x0312

func Listen(user32 *syscall.DLL, keys map[int]*Hotkey, switcher playlistSwitcher, hwnd uintptr) {
	fmt.Println("( ͡° ͜ʖ ͡°) Listening for hot keys in my area...")

	peekmsg := user32.MustFindProc("PeekMessageW")

	var g int // debug goblin
	var debug bool

	for {
		g++
		var msg = &MSG{}
		a, _, _ := peekmsg.Call(uintptr(unsafe.Pointer(msg)), hwnd, WM_HOTKEY, WM_HOTKEY, 1)
		if debug {
			fmt.Printf("%#v %d %d \n", msg, a, g)
		}

		// no message, skip this bish
		if a == 0 {
			time.Sleep(time.Millisecond * 500) // wait for 500 otherwise it commits cringe
			continue
		}

		if key, ok := keys[int(msg.WPARAM)]; ok {
			fmt.Println("Hotkey was pressed:", key)
			switcher.Switch(key.Id)
		}
	}
}
