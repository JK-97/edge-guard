package config

import (
	"jxcore/app/schema"
	"time"
)

var InterSettings = schema.InterSettings{
	HostIP:    "0.0.0.0",
	ChangeLog: "./changelog.json",
	Cmd: schema.Cmd{
		DpkgPR: "dpkg -l | grep ^rc | cut -d' ' -f3 | sudo xargs dpkg --purge",
		DpkgRc: "dpkg -p {} | dpkg -r {}",
	},
	DelPath: []string{
		"/app/init",
		"/app/keys.zip",
		"/app/uuid",
		"/data/database/mongo/adm/data/",
		"/data/database/mongo/adm/oplog.timestamp",
	},
	CleanPath: []string{
		"/app/openvpn_files",
		"/app/platform",
	},
	Restore: schema.Restore{
		DockerPkg: "/restore/dockerimage/",
		DockerDesc: schema.Dockerdesc{
			Name: "desc.json",
			Path: "/restore/dockerimage/desc.json",
		},
		PythonPkg:              "/restore/python_pkg/",
		InstallPythonPkgPrefix: "pip3 install -i http://pypi.jiangxingai.com/simple/ --trusted-host pypi.jiangxingai.com {}",
		MongoPath:              "/restore/apt_pkg/mongo",
	},
	Keys: schema.Keys{
		ZipFile:    "./keys.zip",
		ScriptPath: "/restore/script.sh",
		TmpFile:    "",
		Ceph: schema.Ceph{
			PathInZip: "./keys/etc/ceph/",
			EtcCeph:   "/etc/ceph",
		},
		Openvpn: schema.Openvpn{
			PathInZip:  "./keys/app/openvpn_files/",
			EtcOpenVpn: "/etc/openvpn/client/keys",
			EtcClient:  "/etc/openvpn/client.ovpn",
		},
	},
	Register: schema.Register{
		Host:        "/etc/hosts",
		RegisterKey: "vpnserver.jiangxingai.com\n",
	},
	AutoStart: schema.Autostart{
		Path:  "/etc/profile",
		Shell: "/bin/bash /restore/script.sh &",
	},
	Mongodb: schema.Mongodb{
		Host:              "127.0.0.1",
		MongodbConf:       "/data/database/mongo/adm/base.conf",
		MongoDataPath:     "/data/database/mongo/adm/data/",
		MongodbSupervisor: "adm_mongo",
		MongoPkg: []string{
			"mongodb-org-server",
			" mongodb-org-shell",
			"mongodb-org-tools",
			//"mongodb-org-mongos",
		},
	},

	Redis: schema.Redis{
		Host:         "127.0.0.1",
		UninstallCMD: "apt autoremove -y redis",
	},
	Supervisor: schema.Supervisor{Host: "http://127.0.0.1:8999"},
	FileMonitor: schema.FileMonitor{
		OverSeePath: []string{
			"/app",
			"/edge/...",
		},
		Cleantimestep: 10,
		CleanStrategy: []schema.CleanStrategy{
			{
				Path: "/test",
				Size: 10240,
			},
			{
				Path: "/test1",
				Size: 20000000,
			},
		},
	},
	Gateway: schema.Gateway{
		Name:   "/edge/gateway/bin/gateway",
		Config: "/edge/gateway/bin/gateway.cfg",
	},

	HeartBeat: schema.HearBeat{
		MasterAddr: "10.208.0.1:30431",
		Interval:   time.Duration(3000),
	},
}
