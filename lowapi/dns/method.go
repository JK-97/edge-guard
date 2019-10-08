package dns

import (
    "io"
    "io/ioutil"
    "jxcore/log"
    "math/rand"
    "net"
    "os"
    "os/exec"
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

    // }
    //每一行
    res = strings.Split(string(datastr), "\n")
    for _, rawLine := range res {

        rawLine = strings.TrimSpace(rawLine)
        if string(rawLine[0]) == "#" {
            continue
        }

        if strings.Contains(rawLine, "127.0.0.1") {
            continue
        }

        if pos := strings.Index(rawLine, "nameserver"); pos != -1 {
            server := strings.TrimSpace(rawLine[pos:])
            f.WriteString("server=" + server + "\n")
        }
    }
    
    ResetResolv()
    RestartDnsmasq()
}

func RestartDnsmasq() {
    exec.Command("/bin/bash", "-c", "systemctl restart dnsmasq").Run()
}

func ResetResolv() {
    datatowrite := []byte(hostsRecord)

    log.Info("Write /etc/resolv.conf")
    err := ioutil.WriteFile("/etc/resolv.conf", datatowrite, 0644)
    if err != nil {
        log.Error(err)
    }

}
