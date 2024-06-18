package util

type Keybinding struct {
	Key         string
	Description string
	Action      func()
}
