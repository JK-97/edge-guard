package client

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

// firstword
var rpcCommands = []prompt.Suggest{
	{Text: "status", Description: "Get all program status info"},
	{Text: "tail", Description: "tail the process stderr or stdout"},
	{Text: "stop", Description: "Stop the program"},
	{Text: "start", Description: "Start the program"},
	{Text: "pid", Description: "Get pid by name"},
	{Text: "restart", Description: "Restart program"},
	{Text: "reset", Description: "Restart jxcore"},
	{Text: "log", Description: "jxcore log"},
	// {Text: "signal", Description: "Signal a process"},
	// {Text: "reload", Description: "reload program config"},
}

var logLevel = []prompt.Suggest{
	{Text: "stdout"}, // valid only for federation apiservers
	{Text: "stderr"},
}

var configFile = []prompt.Suggest{
	{Text: "dnsmasqConf", Description: "configFile for dnsmasq"}, // valid only for federation apiservers
	{Text: "initFile", Description: "initFile build by jxcore bootstrap"},
	{Text: "dnsmasqHost", Description: "hostFile for dnsmasq"},
	{Text: "dnsmasqResolv", Description: "resolvFile for dnsmasq"},
}

// 提示参数主逻辑
func (c *customcompleter) argumentsCompleter(args []string, d prompt.Document) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(rpcCommands, args[0], true)
	}

	firstword := args[0]
	switch firstword {
	case "status", "tail", "stop", "start", "pid", "reset", "restart":
		secondword := args[1]
		if len(args) == 2 {
			return prompt.FilterHasPrefix(processTypes, secondword, true)
		}
		if len(args) == 3 {
			switch firstword {
			case "tail":
				thirdWord := d.GetWordBeforeCursor()
				switch thirdWord {
				case "st", "s":
					return logLevel
				case "stde", "stder", "stderr":
					return []prompt.Suggest{logLevel[1]}
				case "stdo", "stdou", "stdout":
					return []prompt.Suggest{logLevel[0]}
				default:
					return logLevel
				}

			}
		}
	case "log":
		cmd := exec.Command("tail", "-f", "/edge/logs/jxcore.log")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
		}

	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}
