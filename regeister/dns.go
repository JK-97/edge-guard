package regeister

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

	"github.com/rjeczalik/notify"
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

const hostsRecord string = "nameserver 127.0.0.1"

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

// DnsFileListener 检测 resolv 文件的改动
func DnsFileListener() {
	ResolvGuard()

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 2)

	// Set up a watchpoint listening for inotify-specific events within a
	// current working directory. Dispatch each InMovedFrom and InMovedTo
	// events separately to c.

	// , notify.InMovedFrom, notify.InMovedTo, notify.Remove
	if err := notify.Watch("/etc/", c, notify.Create, notify.InModify); err != nil {
		log.Error(err)
	}
	defer notify.Stop(c)

	// Inotify reports move filesystem action by sending two events tied with
	// unique cookie value (uint32): one of the events is of InMovedFrom type
	// carrying move source path, while the second one is of InMoveTo type
	// carrying move destination path.

	// Wait for moves.
	for ei := range c {
		switch ei.Event() {
		case notify.InModify, notify.Create:
			if ei.Path() == "/etc/resolv.conf" {
				ResolvGuard()
			}
		}
	}
}
