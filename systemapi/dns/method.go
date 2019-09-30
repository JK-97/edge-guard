package dns

import (
    "io"
    "io/ioutil"
    "jxcore/log"
    "math/rand"
    "net"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

func LookUpDns(domain string) {
    ipRecords, _ := net.LookupIP(domain)
    Shuffle(ipRecords)
    f, err := os.OpenFile("/etc/dnsmasq.d/dnsmasq.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        log.Info(err)
        return
    }
    defer f.Close()
    for _, ip := range ipRecords {
        f.WriteString("server=/.iotedge/" + ip.String() + "\n")
    }

}

// Shuffle 打乱 DNS 记录
func Shuffle(slice []net.IP) {
    r := rand.New(rand.NewSource(time.Now().Unix()))
    for len(slice) > 0 {
        n := len(slice)
        randIndex := r.Intn(n)
        slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
        slice = slice[:n-1]
    }
}

// ResolvGuard 控制 resolv.conf
func ResolvGuard() {
    data, err := ioutil.ReadFile("/etc/resolv.conf")
    if err != nil {
        log.Error(err)
    }
    datastr := string(data)
    if datastr == hostsRecord {
        return
    }
    //每一行
    res := strings.Split(string(datastr), "\n")

    f, err := os.OpenFile("/etc/dnsmasq.d/resolv.conf", os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Error("Open /etc/dnsmasq.d/resolv.conf", err)
    } else {
        f.Seek(0, io.SeekStart)
    }
    defer f.Close()

    // var resolvbak = make([]byte, 0)
    // f.Read(resolvbak)
    // ioutil.WriteFile("/etc/dnsmasqs.d/resolv.bak", resolvbak, 0644)
    // if err != nil {

    // }
    for _, percol := range res {

        percol = strings.TrimSpace(percol)
        if len(percol) > 19 {
            if string(percol[0]) == "#" {
                log.Error("no start with # ")
                //每一句不以#开头
            } else {
                r := regexp.MustCompile(" {0,}nameserver {1,}127.0.0.1")
                rip := regexp.MustCompile(" {0,}nameserver {1,}((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))")
                if r.MatchString(percol) {
                    //匹配nemeserver 127.0.0.1
                } else {
                    res := rip.FindAllStringSubmatch(percol, -1)
                    for _, perip := range res {
                        ip := perip[1]
                        f.WriteString("server=" + ip + "\n")
                    }

                }

            }
        }

    }

    //删除resolve

    datatowrite := []byte(hostsRecord)

    log.Info("Write /etc/resolv.conf")
    err = ioutil.WriteFile("/etc/resolv.conf", datatowrite, 0644)
    if err != nil {
        log.Error(err)
    }

    exec.Command("/bin/bash", "-c", "systemctl restart dnsmasq").Run()
}
