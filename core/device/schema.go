package device

type Device struct {
    WorkerID     string `yaml:"workerid"`
    Key        string `yaml:"key"`
    DhcpServer string `yaml:"dhcpserver"`
    Vpn        string `yaml:"vpn"`
}

type buildkeyreq struct {
    Workerid string `json:"wid"`
    Ticket   string `json:"ticket"`
}
type data struct {
    Key         string `json:"key"`
    DeadLine    string `json:"deadLine"`
    RemainCount string `json:"remainCount"`
}

type buildkeyresp struct {
    Data data `json:"data"`
}
