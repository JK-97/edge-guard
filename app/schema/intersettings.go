package schema

import (
	"time"
)

type InterSettings struct {
	HostIP    string   `json:"host_ip"`
	Cmd       Cmd      `json:"cmd"`
	DelPath   []string `json:"del_path"`
	CleanPath []string `json:"clean_path"`
	Restore   Restore  `json:"restore"`
	Keys      Keys     `json:"Keys"`
	ChangeLog string   `json:"change_log"`

	Register    Register    `json:"register"`
	AutoStart   Autostart   `json:"auto_start"`
	Mongodb     Mongodb     `json:"mongodb"`
	Redis       Redis       `json:"redis_host"`
	Supervisor  Supervisor  `json:"supervisor"`
	FileMonitor FileMonitor `json"filemonitor"`
	Gateway     Gateway     `json:"gateway"`

	HeartBeat   HearBeat    `json:"heart_beat"`
}
type HearBeat struct {
	Interval time.Duration
	MasterAddr string
}
type CleanStrategy struct {
	Path string `json:"path"`
	Size int    `json:"size"`
}


type Dockerdesc struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
type Restore struct {
	DockerPkg              string     `json:"docker_pkg"`
	DockerDesc             Dockerdesc `json:"docker_desc"`
	PythonPkg              string     `json:"python_pkg"`
	InstallPythonPkgPrefix string     `json:"install_python_pkg_prefix"`
	MongoPath              string     `json:"mongo_path"`
}

type Cmd struct {
	DpkgRc string `json:"dpkg_rc"`
	DpkgPR string `json:"dpkg_p_r"`
}

type Autostart struct {
	Path  string `json:"path"`
	Shell string `json:"shell"`
}

type Gateway struct {
	Name   string `json:"name"`
	Config string `json:"Config"`
}

type Mongodb struct {
	Host              string   `json:"host"`
	MongoDataPath     string   `json:"mongodatapath"`
	MongodbConf       string   `json:"mongodb_conf"`
	MongodbSupervisor string   `json:"mongodb_supervisor"`
	MongoPkg          []string `json:"mongo_pkg"`
}

type Redis struct {
	Host         string `json:"redis_host"`
	UninstallCMD string `json:"uninstall_redis_cmd"`
}
type Supervisor struct {
	Host string `json:"supervisor_host"`
}

type FileMonitor struct {
	OverSeePath   []string        `json:"oversee_path"`
	Cleantimestep int             `json:"cleantimestep"`
	CleanStrategy []CleanStrategy `json:"clean_strategy"`
}
type Ceph struct {
	PathInZip string `json:"pathInzip"`
	EtcCeph   string `json:"EtcCeph"`
}
type Openvpn struct {
	PathInZip  string `json:"pathInzip"`
	EtcOpenVpn string `json:"EtcOpenVpn "`
	EtcClient  string `json:"EtcClient"`
}

type Keys struct {
	ZipFile string `json:"zipfile"`

	ScriptPath string  `json:"script_path""`
	TmpFile    string  `json:"tmpfile"`
	Ceph       Ceph    `json:"ceph"`
	Openvpn    Openvpn `json:"openvpn"`
}

type Register struct {
	Host        string `json:"host"`
	RegisterKey string `json:"registerKey"`
}
