package serve

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"jxcore/gateway/utils"
)

const nodeInfoPath = "/edge/init"

// HandleGetNodeInfo 获取节点的信息
func HandleGetNodeInfo(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(nodeInfoPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	resp := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := strings.SplitN(scanner.Text(), ":", 2)
		l0 := strings.TrimSpace(l[0])
		if len(l) == 2 {
			resp[l0] = strings.TrimSpace(l[1])
		} else {
			resp[l0] = ""
		}
	}
	utils.WriteJson(w, resp)
}
