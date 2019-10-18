package dns


type consulConfig struct {
    Server           bool     `json:"server"`
    ClientAddr       string   `json:"client_addr"`
    AdvertiseAddrWan string   `json:"advertise_addr_wan"`
    BootstrapExpect  int      `json:"bootstrap_expect"`
    Datacenter       string   `json:"datacenter"`
    NodeName         string   `json:"node_name"`
    RetryJoinWan     []string `json:"retry_join_wan"`
    UI               bool     `json:"ui"`
}
