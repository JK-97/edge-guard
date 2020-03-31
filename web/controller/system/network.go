package system

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/JK-97/edge-guard/internal/network/dns"
	"github.com/JK-97/edge-guard/internal/network/iface"
	"github.com/JK-97/edge-guard/lowapi/system"
	"github.com/JK-97/edge-guard/oplog"
	"github.com/JK-97/edge-guard/oplog/logs"
	"github.com/JK-97/edge-guard/oplog/types"
	"github.com/JK-97/edge-guard/web/controller/utils"

	"github.com/JK-97/go-utils/logger"
	"github.com/gorilla/mux"
)

type interfaceInfo struct {
	Name       string   `json:"name"`
	Auto       bool     `json:"auto"`
	IP         string   `json:"ip"`
	Mask       int      `json:"mask"`
	Gateway    string   `json:"gateway"`
	Nameserver []string `json:"nameserver"`
}

func GetNetworkInterfaces(w http.ResponseWriter, r *http.Request) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var resp struct {
		Interfaces []interfaceInfo `json:"interfaces"`
	}
	for _, i := range interfaces {
		if strings.HasPrefix(i.Name, "veth") ||
			strings.HasPrefix(i.Name, "br-") ||
			strings.HasPrefix(i.Name, "docker") ||
			strings.HasPrefix(i.Name, "lo") {
			continue
		}

		info, err := parseInterfaceInfo(&i)
		if err != nil {
			logger.Error(err)
			continue
		}
		resp.Interfaces = append(resp.Interfaces, info)
	}

	utils.RespondSuccessJSON(resp, w)
}

func parseInterfaceInfo(i *net.Interface) (interfaceInfo, error) {
	addrs, err := i.Addrs()
	if err != nil {
		return interfaceInfo{}, err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		default:
			continue
		}
		if ip.To4() == nil {
			continue
		}

		maskOnes, _ := ip.DefaultMask().Size()
		info := interfaceInfo{
			Name: i.Name,
			IP:   ip.String(),
			Mask: maskOnes,
		}

		if gwRoute, _ := iface.GetGWRoute(i.Name); gwRoute != nil {
			info.Gateway = gwRoute.Gw.String()
		}

		if nameservers, err := dns.ParseInterfaceDNSResolv(i.Name); err == nil {
			info.Nameserver = nameservers
		}

		return info, nil
	}
	return interfaceInfo{}, fmt.Errorf("network info for interface \"%v\" not found", i.Name)
}

func GetNetworkInterfaceByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ifaceName := vars["iface"]
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		panic(err)
	}

	info, err := parseInterfaceInfo(iface)
	if err != nil {
		panic(err)
	}
	utils.RespondSuccessJSON(info, w)
}

type fourGInfo struct {
	Enable         bool    `json:"enable"`
	SignalStrength float64 `json:"signalstrength"`
	IP             string  `json:"ip"`
}

var fourGInterface = "usb0"

func GetFourGInterface(w http.ResponseWriter, r *http.Request) {
	scpritPath := "/usb/locl/sbin/G8100_NoMCU.py"

	resp := &fourGInfo{}

	iface, err := net.InterfaceByName(fourGInterface)
	if err != nil {
		resp.Enable = false
		utils.RespondSuccessJSON(resp, w)
		return
	}
	ifaceInfo, _ := parseInterfaceInfo(iface)

	output, err := system.RunCommandWithOutput(fmt.Sprintf("python %s CMD AT+CSQ", scpritPath))
	if err != nil {
		utils.RespondReasonJSON(resp, w, fmt.Sprintf("exec scprit %s with err,%s", scpritPath, err.Error()), 400)
		return
	}

	signalStrengthRssi, err := parseFourGInfo(output)
	if err != nil {
		utils.RespondReasonJSON(resp, w, fmt.Sprintf("prase scprit %s output with err,%s", scpritPath, err.Error()), 400)
		return
	}

	resp.Enable = true
	resp.IP = ifaceInfo.IP
	resp.SignalStrength = 2*signalStrengthRssi - 113

	utils.RespondSuccessJSON(resp, w)

}

func parseFourGInfo(output []byte) (float64, error) {
	prefix := "+CSQ:"
	res := ""
	data := strings.Split(string(output), "\n")
	for _, line := range data {
		if strings.HasPrefix(line, prefix) {
			res = strings.TrimSpace(strings.Trim(line, prefix))
		}
	}
	res = strings.ReplaceAll(res, ",", ".")
	logger.Info("csq:" + res)
	return strconv.ParseFloat(res, 64)

}

func EnableFourGInterface(w http.ResponseWriter, r *http.Request) {
	queryValue := r.URL.Query()
	ifenable := queryValue.Get("enable")
	enable, err := strconv.ParseBool(ifenable)
	if err != nil {
		utils.RespondReasonJSON(nil, w, "pare args error", 400)
		return
	}

	if enable {
		err = system.RunCommand(fmt.Sprintf("ifup %s", fourGInterface))
	} else {
		err = system.RunCommand(fmt.Sprintf("ifdown %s", fourGInterface))
	}
	if err != nil {
		utils.RespondReasonJSON(nil, w, fmt.Sprintf("operated faild with err : %s", err.Error()), 400)
		return
	}
	oplog.Insert(logs.NewOplog(types.NETWORKE, fmt.Sprintf("set 4g %s", enable)))
	utils.RespondSuccessJSON(nil, w)
}
