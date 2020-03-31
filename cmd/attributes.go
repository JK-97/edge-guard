package cmd

const (
	LogsPath           = "/edge/logs"
	TargetVersionFile  = "/etc/edgetarget"
	CurrentVersionFile = "/edge/VERSION"
)

var DependsImages = []string{}

var BinFilesMAP = map[string]string{}

var DependsFile = []string{}
