package network

import (
    "fmt"
    "jxcore/core/device"
    "jxcore/log"
    "jxcore/lowapi/utils"
    "jxcore/lowapi/vpn"
    "net"
    "os/exec"
    "strings"
    "time"
)

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
        log.WithFields(log.Fields{"Operating": "GetEthIP"}).Error(err)
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



// GetClusterIP 获取集群内网 VPN IP
func GetClusterIP() string {
    d,err:= device.GetDevice()
    utils.CheckErr(err)
    switch d.Vpn {
    case device.VPNModeOPENVPN:
        tun0interface, err := GetMyIP(vpn.OpenVPNInterface)
        if err != nil {
            log.WithFields(log.Fields{"Operating": "GetClusterIP"}).Error(err)
            return ""
        }
        return tun0interface
    case device.VPNModeWG:
        wg0interface, err := GetMyIP(vpn.WireGuardInterface)
        if err != nil {
            log.WithFields(log.Fields{"Operating": "GetClusterIP"}).Error(err)
            return ""
        }
        return wg0interface
    }
    return ""
}


func CheckNetwork() bool {
    cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
    err := cmd.Run()
    if err != nil {
        fmt.Println(err.Error())
        return false
    } else {
        log.WithFields(log.Fields{"Device":"NetWork Check"}).Info("Net Status , OK")
    }
    return true

}