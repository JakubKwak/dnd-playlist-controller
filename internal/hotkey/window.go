package hotkey

import (
	"fmt"
	"syscall"
	"unsafe"
)

// todo get rid >:(
var (
	user32             = syscall.NewLazyDLL("user32.dll")
	procCreateWindowEx = user32.NewProc("CreateWindowExW")
)

func createWindow(
	dwExStyle uint32,
	lpClassName, lpWindowName *uint16,
	dwStyle uint32,
	x, y, nWidth, nHeight int,
	hWndParent, hMenu, hInstance uintptr,
	lpParam unsafe.Pointer,
) uintptr {
	ret, _, _ := procCreateWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(unsafe.Pointer(lpClassName)),
		uintptr(unsafe.Pointer(lpWindowName)),
		uintptr(dwStyle),
		uintptr(x),
		uintptr(y),
		uintptr(nWidth),
		uintptr(nHeight),
		hWndParent,
		hMenu,
		hInstance,
		uintptr(lpParam),
	)
	return ret
}

// create a simple window so we can listen to messages only for this app and not all messages on this thread
func GiveSimpleWindowPls() uintptr {
	var hwnd uintptr
	className, _ := syscall.UTF16PtrFromString("STATIC")
	windowName, _ := syscall.UTF16PtrFromString("Sample Window")

	hwnd = createWindow(
		0,
		className,
		windowName,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		nil,
	)

	if hwnd != 0 {
		fmt.Println("HWND:", hwnd)
	} else {
		fmt.Println("window fricking gone man")
	}
	return hwnd
}
