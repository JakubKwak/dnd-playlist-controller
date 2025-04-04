package hotkey

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

var (
	modsMap = map[string]int{
		"alt":   ModAlt,
		"ctrl":  ModCtrl,
		"shift": ModShift,
		"win":   ModWin,
	}
)

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

func ParseMod(str string) (int, bool) {
	c, ok := modsMap[strings.ToLower(str)]
	return c, ok
}

// returns "Hotkey[Id: 1, Alt+Ctrl+O]" for logging
func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}
