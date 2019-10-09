package updatemange

type UpgradeProcess struct {
    Status     UpgradeStatus `json:"status"`
    Target     VersionInfo   `json:"target"`
    NowVersion VersionInfo   `json:"now_version"`
}

type targetversionfile struct {
    Target map[string]string `json:"target"`
}

type Reqdatastruct struct {
    Data map[string]string `json:"data"`
}

type Respdatastruct struct {
    Status   string            `json:"status"`
    WorkerId string            `json:"worker_id"`
    PkgInfo  map[string]string `json:"pkg_info"`
}
type ComponentInfo struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}
