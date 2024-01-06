package commands

type Prompt struct {
	Items   []string
	Handler func(int, string) bool
}
