package device

type Device struct {
    WorkID     string `json:"workerid"`
    Key        string `json:"key"`
    DhcpServer string `json:"dhcpserver"`
    Vpn        string `json:"vpn"`
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
