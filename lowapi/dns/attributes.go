package dns

const (
    HostsFile       string = "/etc/dnsmasq.hosts"
    hostsRecord     string = "nameserver 127.0.0.1"
    LocalHostName   string = "edgegw.localhost"
    IotedgeHostName string = "edgegw.iotedge"
    MasterHostName  string = "master.iotedge"
    ResolvFile      string = "/etc/resolv.conf"

    consulConfigPath string = "/data/edgex/consul/config/consul_conf.json"
)
