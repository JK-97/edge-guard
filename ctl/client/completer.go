package client

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

type customcompleter struct {
}

// NewCompleter 返customcompleter 对象
func NewCompleter() *customcompleter {
	return &customcompleter{}
}

// 用于补全的接口
func (c *customcompleter) completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return rpcCommands
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	return c.argumentsCompleter(args, d)
}
