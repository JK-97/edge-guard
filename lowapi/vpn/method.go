package vpn

import (
    "bufio"
    "jxcore/log"
    "os"
    "os/exec"
    "strings"
    "sync"
    "time"
)

// StartWg 打开 WireGuard VPN
func StartWg() error {
    cmd := exec.Command("wg-quick", "up", WireGuardInterface)
    out, err := cmd.CombinedOutput()
    if err == nil {
        log.Info("wg up success")
    } else {
        log.Info("wg up failed: ", err, string(out))
    }

    return err
}

// CloseWg 关闭 WireGuard VPN
func CloseWg() error {
    cmd := exec.Command("wg-quick", "down", WireGuardInterface)
    out, err := cmd.Output()
    if err == nil {
        log.Info("wg down success")
    } else {
        log.Info("wg down failed: ", err, string(out))
    }
    return err
}



// Startopenvpn 打开 OpenVPN
func Startopenvpn() error {
    cmd := exec.Command("openvpn", openvpnConfigPath)
    pipe, err := cmd.StdoutPipe()
    err = cmd.Start()
    if err != nil {
        log.Info("openvpn up failed :", err)
        return err
    }
    scanner := bufio.NewScanner(pipe)
    // 检测 OpenVPN 是否正常启动
    wg := new(sync.WaitGroup)
    wg.Add(1)
    go func() {
        for scanner.Scan() {
            if strings.Contains(scanner.Text(), openvpnSuccessMessage) {
                pipe.Close()
                wg.Done()
                break
            }
        }
        if scanner.Err() == nil {
            return
        }
        if perr, ok := scanner.Err().(*os.PathError); ok {
            if perr.Err != os.ErrClosed {
                err = scanner.Err()
            }
        } else {
            err = scanner.Err()
        }
    }()
    timer := time.NewTimer(openvpnWaitTimeout)
    go func() {
        select {
        case <-timer.C:
            err = ErrOpenVPNTimeout
            pipe.Close()
            return
        }
    }()
    wg.Wait()
    if err == nil {
        log.Info("openvpn up success")
    }
    return err
}


// Closeopenvpn 关闭 OpenVPN
func Closeopenvpn() error {
    c := "killall openvpn"
    cmd := exec.Command("/bin/sh", "-c", c)
    err := cmd.Run()
    if err == nil {
        log.Info("openvpn down success")
    } else {
        log.Info("openvpn down failed ", err)
    }
    return err
}
