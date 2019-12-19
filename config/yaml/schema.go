package yaml

type YamlSchema struct {
	Cri              bool `yaml:"cri"`
	Jxserving        bool `yaml:"jxserving"`
	FsIndex          bool `yaml:"fsindex"`
	Speaker          bool `yaml:"speaker"`
	Devicemanagement struct {
		Camera    bool `yaml:"camera"`
		Rs485     bool `yaml:"rs485"`
		Aiserving bool `yaml:"aiserving"`
	} `yaml:"devicemanagement"`
	Monitor struct {
		Telegraf bool `yaml:"telegraf"`
	} `yaml:"monitor"`
	Tools struct {
		Nettools struct {
			Ifplugd bool `yaml:"ifplugd"`
		} `yaml:"nettools"`
		Mcutools struct {
			Watchdog        bool `yaml:"watchdog"`
			Powermanagement bool `yaml:"powermanagement"`
			Mcuserver       bool `yaml:"mcuserver"`
		} `yaml:"mcutools"`
		Nodetools struct {
			Cleaner      bool `yaml:"cleaner"`
			Usblistener  bool `yaml:"usblistener"`
			Filelistener bool `yaml:"filelistener"`
		} `yaml:"nodetools"`
	} `yaml:"tools"`
	Synctools struct {
		Vpn    bool `yaml:"vpn"`
		Db     bool `yaml:"db"`
		Tsdb   bool `yaml:"tsdb"`
		Mq     bool `yaml:"mq"`
		Fs     bool `yaml:"fs"`
		Config bool `yaml:"config"`
	} `yaml:"synctools"`
	FixedResolver string `yaml:"fixedresolver"`
	Debug         bool   `yaml:"debug"`

	IFace struct {
		Priority       []string `yaml:"priority"`
		Backup         string   `yaml:"backup"`
		SwitchInterval string   `yaml:"switch_interval"`
	} `yaml:"iface"`

	MountCfg map[string]string `yaml:"mount_cfg"`
}
