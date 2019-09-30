package dnsdetector

import (
    "github.com/rjeczalik/notify"
    "jxcore/log"
    "jxcore/systemapi/dns"
)

// DnsFileListener 检测 resolv 文件的改动
func DnsFileListener() {
    dns.ResolvGuard()
    
    c := make(chan notify.EventInfo, 2)
    if err := notify.Watch("/etc/", c, notify.Create, notify.InModify); err != nil {
        log.Error(err)
    }
    defer notify.Stop(c)
    
    
    for ei := range c {
        switch ei.Event() {
        case notify.InModify, notify.Create:
            if ei.Path() == "/etc/resolv.conf" {
                dns.ResolvGuard()
            }
        }
    }
}
