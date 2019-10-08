package register

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "io/ioutil"
    "jxcore/core/device"
    "jxcore/log"
    "jxcore/lowapi/dns"
    "jxcore/lowapi/utils"
    "jxcore/lowapi/vpn"
    "jxcore/version"
    "net/http"
    "os/exec"
    "time"
)

// findMasterFromDHCPServer 从 DHCP 服务器 获取 Master 节点的 IP
func findMasterFromDHCPServer(workerid string, key string) (masterip string, err error) {
    
    currentdevice,err:=device.GetDevice()
    utils.CheckErr(err)
    
    
    reqinfo := reqRegister{
        WorkerID: workerid,
        Nonce:    time.Now().Unix(),
        Key:      key,
        Version:  version.Version,
    }

    reqdata, err := json.Marshal(reqinfo)
    if err != nil {
        log.Error(err)
    }
    //req base64加密
    n := enc.EncodedLen(len(reqdata))
    dst := make([]byte, n)
    enc.Encode(dst, reqdata)

    //通过dhcpserver获取key
    reqbody := bytes.NewBuffer(dst)

    url := currentdevice.DhcpServer + wireguardRegisterPath
    if currentdevice.Vpn == device.VPNModeOPENVPN {
        url = currentdevice.DhcpServer + openvpnRegisterPath
    }

    resp, err := http.Post(url, "application/json", reqbody)
    if err != nil {
        log.Error(err, "restart dnsmasq")
        exec.Command("service", "dnsmasq", "restart").Run()
        getmymaster(workerid, key)
        return
    } else if resp.StatusCode != http.StatusOK {
        err = errors.New(resp.Status)
        return
    }

    masterip = resp.Header.Get("X-Master-IP")
    defer resp.Body.Close()

    if masterip != "" {
        ip := dns.FindMasterFromHostFile()

        dns.UpdateMasterIPToHosts(masterip)
        if masterip != ip {
            dns.OnMasterIPChanged(masterip)
        }
    }

    //获得加密wgkey zip
    buff, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Error(err)
        return
    }

    //解密
    r := ChaosReader{
        Bytes:  buff,
        Offset: prefix,
    }
    content := make([]byte, len(buff)-prefix-suffix)
    _, err = r.Read(content)

    log.Info("Updating VPN")
    // 替换vpn配置
    switch currentdevice.Vpn {
    case device.VPNModeWG:
        log.Info("VPN Mode: ", currentdevice.Vpn)
        replacesetting(bytes.NewReader(content), "/etc/wireguard")
        //vpn commponet检测配置变动 启动wireguard ,wg0
        vpn.CloseWg()
        if err :=vpn. StartWg(); err == nil {
            dns.OnVPNConnetced()
        }

    case device.VPNModeOPENVPN:
        log.Info("VPN Mode: ", currentdevice.Vpn)
        replacesetting(bytes.NewReader(content), "/etc/openvpn/")
        vpn.Closeopenvpn()
        if err := vpn.Startopenvpn(); err == nil {
            dns.OnVPNConnetced()
        }
    default:
        log.Error("err model")
        return
    }

    exec.Command("service", "dnsmasq", "restart").Run()
    return
}


func (r *ChaosReader) Read(p []byte) (n int, err error) {
    length := len(r.Bytes)
    remain := length - r.Offset
    if remain <= 0 {
        return 0, io.EOF
    }
    length = len(p)
    if length > remain {
        err = io.EOF
    } else {
        remain = length
    }

    for n = 0; n < remain; n++ {
        b := r.Bytes[r.Offset+n]
        if b >= 0x80 {
            p[n] = b - 0x80
        } else {
            p[n] = b + 0x80
        }
    }
    r.Offset += remain
    return
}


func replacesetting(formfile *bytes.Reader, toetc string) {
    formfile.Seek(0, io.SeekStart)

    buff, err := ioutil.ReadAll(formfile)
    if err != nil {
        log.Error(err)
    }

    //time.Sleep(2 * time.Second)
    err = utils.Unzip(buff, toetc)
    if err != nil {
        log.Error(err)
    }
}


func getmymaster(workerid, key string) (mymasterip string, err error) {
    masterip := dns.FindMasterFromHostFile()

    if masterip == "" {
        masterip, err = findMasterFromDHCPServer(workerid, key)
    }
    if err != nil {
        return
    }
    dns.AppendHostnameHosts(workerid)

    log.Info("Finish Update VPN")
    // _, errusb0 := GetMyIP("usb0")

    return masterip, err
}
