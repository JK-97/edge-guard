package system

import (
	"github.com/JK-97/edge-guard/lowapi/system"
	"github.com/JK-97/edge-guard/monitor"
	"github.com/JK-97/edge-guard/web/controller/utils"
	"net/http"
	"strconv"
	"strings"
)

type tfCardInfo struct {
	Mounted   bool `json:"mounted"`
	All       int  `json:"all"`
	Used      int  `json:"used"`
	Available int  `json:"available"`
}

func GetStorgeInfo(w http.ResponseWriter, r *http.Request) {
	responce := &tfCardInfo{}
	ifMounted, err := monitor.CheckMount("/media/mmcblk1p1")
	responce.Mounted = ifMounted
	if err != nil {
		return
	}

	if !responce.Mounted {
		utils.RespondSuccessJSON(responce, w)
		return
	}

	output, err := system.RunCommandWithOutput("df /media/mmcblk1p1")
	if err != nil {
		utils.RespondReasonJSON(nil, w, "run command failed", 400)
		return
	}

	if strings.Contains(string(output), "No such file or directory") {
		utils.RespondReasonJSON(nil, w, "No such file or directory", 400)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	data := strings.Fields(lines[1])
	all, _ := strconv.Atoi(data[1])
	used, _ := strconv.Atoi(data[2])
	available, _ := strconv.Atoi(data[3])

	responce.All = all
	responce.Used = used
	responce.Available = available
	utils.RespondSuccessJSON(responce, w)
}
