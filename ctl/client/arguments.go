package client

import (
	"github.com/c-bata/go-prompt"
)

// firstword
var rpcCommands = []prompt.Suggest{
	{Text: "status", Description: "Get all program status info"},
	// {Text: "tail", Description: "tail the process stderr or stdout"},
	{Text: "stop", Description: "Stop the program"},
	{Text: "start", Description: "Start the program"},
	{Text: "pid", Description: "Get pid by name"},
	{Text: "restart", Description: "Restart program"},
	{Text: "reset", Description: "Restart jxcore"},
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
func (c *customcompleter) argumentsCompleter(args []string) []prompt.Suggest {
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
				return logLevel
			}
		}

	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}
