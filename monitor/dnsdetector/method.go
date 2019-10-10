package dnsdetector

import (
    "github.com/rjeczalik/notify"
    "jxcore/log"
    "jxcore/lowapi/dns"
    "os"
)

// DnsDetector 检测 resolv 文件的改动
func DnsDetector() {

    c := make(chan notify.EventInfo, 2)
    if err := notify.Watch(resolvfile, c, notify.All); err != nil {
        log.Error(err)
    }
    for ei := range c {
        switch ei.Event() {
        case notify.Remove:
            dns.ResolvGuard()
            notify.Stop(c)
            DnsDetector()

        }
    }
}
func RunDnsDetector() {
    dns.RestartDnsmasq()
    // TODO check /etc/resolv.conf exists
    if _, err := os.Stat(resolvfile); err == nil {
        dns.ResolvGuard()
    } else {
        log.Info("Has no detect the resolv.conf")
        dns.ResetResolv()
    }
    DnsDetector()
}
