package regeister

import (
	"fmt"
	"jxcore/config"
	"jxcore/log"
	"net"
	"strings"
	"time"
)

// AliveReport 上报心跳
func AliveReport(masterip string) {
	var msg string

	if DeviceKey.Vpn == VPNModeWG {
		time.Sleep(2 * time.Second)
		wg0interface, err := GetMyIP(WireGuardInterface)
		if err != nil {
			log.Error(err)
		}
		msg = fmt.Sprintf(wg0interface + ":" + DeviceKey.WorkID)
	} else if DeviceKey.Vpn == VPNModeOPENVPN {
		time.Sleep(2 * time.Second)
		tun0interface, err := GetMyIP(OpenVPNInterface)
		if err != nil {
			log.Error(err)
		}
		msg = fmt.Sprintf(tun0interface + ":" + DeviceKey.WorkID)
	}
	// TODO: 心跳协议
	// msg = msg+ "\n\n"
	//发送tcp
	conn, err := net.Dial("tcp", masterip+":30431")
	if err != nil {
		log.WithFields(log.Fields{"Operating": "AliveReport"}).Error("disconnect my master", err)
		Patternmatching()

	}
	log.Info("connected")
	Connectable <- true
	defer conn.Close()
	ticker := time.NewTicker(time.Millisecond * config.InterSettings.HeartBeat.Interval)

	select {
	case <-ticker.C:

		for index := 0; index < 5; index++ {
			for range ticker.C {
				_, err := conn.Write([]byte(msg))
				if err != nil {
					log.WithFields(log.Fields{"Operating": "AliveReport"}).Error("disconnect my master", err)
					//心跳断联,获取新master
					time.Sleep(3 * time.Second)
					ticker.Stop()
					break
				}
			}

		}
		log.WithFields(log.Fields{"Operating": "AliveReport"}).Error(" will get a new master in 5 second ", err)
	}
	time.Sleep(5 * time.Second)
	Patternmatching()
}

// GetClusterIP 获取集群内网 VPN IP
func GetClusterIP() string {
	switch DeviceKey.Vpn {
	case VPNModeOPENVPN:
		tun0interface, err := GetMyIP(OpenVPNInterface)
		if err != nil {
			log.WithFields(log.Fields{"Operating": "GetClusterIP"}).Error(err)
			return ""
		}
		return tun0interface
	case VPNModeWG:
		wg0interface, err := GetMyIP(WireGuardInterface)
		if err != nil {
			log.WithFields(log.Fields{"Operating": "GetClusterIP"}).Error(err)
			return ""
		}
		return wg0interface
	}
	return ""
}

// GetMyIP 获取 IP
func GetMyIP(name string) (string, error) {
	a, err := net.InterfaceByName(name)
	if err != nil {

		log.WithFields(log.Fields{"Operating": "GetMyIP"}).Error(err, "  "+name)
		return "", err
	}

	addrs, err := a.Addrs()

	if err != nil {
		log.Error(err)
		return "127.0.0.1", err
	}
	//else {
	//	res := strings.Split(addrs[0].String(), "/")
	//	resstr := strings.ReplaceAll(res[0], ".", "_")
	//	resstr = "worker-"+resstr
	//	exec.Command("hostnamectl", "set-hostname", resstr).Start()
	//
	//}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			return ipNet.IP.String(), err
		}
	}
	return "127.0.0.1", err
}

// GetMyIPSlice 获取指定网卡上的 IP 列表
func GetMyIPSlice(name string) []string {

	a, err := net.InterfaceByName(name)
	if err != nil {

		log.WithFields(log.Fields{"Operating": "GetMyIPSlice"}).Error(err, "  "+name)
		return nil
	}

	addrs, err := a.Addrs()

	if err != nil {
		log.Error(err)
		return nil
	}
	//else {
	//	res := strings.Split(addrs[0].String(), "/")
	//	resstr := strings.ReplaceAll(res[0], ".", "_")
	//	resstr = "worker-"+resstr
	//	exec.Command("hostnamectl", "set-hostname", resstr).Start()
	//
	//}
	result := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			result = append(result, ipNet.IP.String())
		}
	}
	return result
}

// GetEthIP 获取以太网 ip
func GetEthIP() string {
	conn, err := net.Dial("tcp", "114.114.114.114:53")

	for err != nil {
		log.WithFields(log.Fields{"Operating": "GetMyIPSlice"}).Error(err)
		time.Sleep(500 * time.Millisecond)
		conn, err = net.Dial("tcp", "114.114.114.114:53")
	}

	defer conn.Close()
	addr := conn.LocalAddr()
	switch addr.(type) {
	case *net.TCPAddr:
		ip := addr.(*net.TCPAddr)
		return ip.IP.String()
	case *net.IPAddr:
		ip := addr.(*net.IPAddr)
		return ip.IP.String()
	}

	addrStr := conn.LocalAddr().String()
	if strings.Contains(addrStr, "[") {
		addrStr = strings.Split(addrStr, "]")[0]
		return addrStr[1:]
	} else if strings.Contains(addrStr, ":") {
		if h, _, err := net.SplitHostPort(addrStr); err == nil {
			return h
		}
	}

	return addrStr
}
