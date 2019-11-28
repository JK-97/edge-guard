package client

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
)

var processTypes = []prompt.Suggest{}

//启动client
func Run(version string) {
	c := NewCompleter()
	e := NewRpcExcutior(
		"http://127.0.0.1:9001",
		"root",
		"",
	)
	rpcc := e.createRpcClient()
	defer fmt.Println("Bye!")

	e.status(rpcc, nil)
	processNameList := e.getAllProcessesName(rpcc)
	for _, processName := range processNameList {
		processTypes = append(processTypes, prompt.Suggest{Text: processName})
	}

	p := prompt.New(
		e.Execute,
		c.completer,
		prompt.OptionTitle("jxcore-prompt: interactive jxcore client"),
		prompt.OptionPrefix("jxcorectl > "),
		prompt.OptionInputTextColor(prompt.Yellow),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}
