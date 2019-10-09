package device

import (
    "crypto/md5"
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "jxcore/log"
    "jxcore/lowapi/utils"
    "jxcore/version"
    "math/rand"
    "runtime"
    "time"
)

func GetDeviceType() (devicetype string) {
    return version.Type
}

func GetDevice() (device *Device, err error) {
    readdata, err := ioutil.ReadFile("/edge/init")
    utils.CheckErr(err)
    err = yaml.Unmarshal(readdata, &device)
    utils.CheckErr(err)
    return
}

// BuildWokerID 生成wokerid 
func BuildWokerID() string {
    perfilx := "J"
    if runtime.GOARCH == "amd64" {
        perfilx = perfilx + "02"
    } else {
        perfilx = perfilx + "01"
    }
    content, err := ioutil.ReadFile("/proc/cpuinfo")
    utils.CheckErr(err)
    md5info := md5.Sum(content[len(content)-17:])
    md5str := fmt.Sprintf("%x", md5info)
    workerid := perfilx + md5str[len(md5str)-7:]
    return workerid
}

// BuildDeviceInfo
func (d *Device) BuildDeviceInfo(vpnmodel string, ticket string, authHost string) {
    if d.WorkerID == "" {
        d.WorkerID = BuildWokerID()
    }
    if vpnmodel == VPNModeRandom {
        r := rand.New(rand.NewSource(time.Now().Unix()))
        vpnmodel = vpnSlice[r.Intn(len(vpnSlice))]
    }
    if GetDeviceType() == version.Pro {
        //pro
    } else {
        //base
        if vpnmodel != VPNModeLocal || authHost != VPNModeLocal{
            log.Fatal("Base version can not support networking")
        }
        
    }
    d.Vpn = vpnmodel
    d.DhcpServer = authHost
    log.Info("Update Init Config File")
    outputdata, err := yaml.Marshal(d)
    utils.CheckErr(err)
    ioutil.WriteFile("/edge/init", outputdata, 0666)

}
