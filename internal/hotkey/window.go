package hotkey

import (
	"fmt"
	"syscall"
	"unsafe"
)

func createWindow(
	user32 *syscall.DLL,
	dwExStyle uint32,
	lpClassName, lpWindowName *uint16,
	dwStyle uint32,
	x, y, nWidth, nHeight int,
	hWndParent, hMenu, hInstance uintptr,
	lpParam unsafe.Pointer,
) uintptr {
	procCreateWindowEx := user32.MustFindProc("CreateWindowExW")

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
func GiveSimpleWindowPls(user32 *syscall.DLL) (uintptr, error) {
	var hwnd uintptr
	className, _ := syscall.UTF16PtrFromString("STATIC")
	windowName, _ := syscall.UTF16PtrFromString("Simple Window")

	hwnd = createWindow(
		user32,
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
		return 0, fmt.Errorf("window fricking gone man")
	}
	return hwnd, nil
}
