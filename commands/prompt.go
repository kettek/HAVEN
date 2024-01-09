package commands

type Prompt struct {
	Items   []string
	Message string
	Handler func(int, string) bool
}
