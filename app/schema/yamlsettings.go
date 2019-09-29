package schema

type YamlSchema struct {
	Cri              bool `json:"cri"`
	Devicemanagement struct {
		Camera    bool `json:"camera"`
		Rs485     bool `json:"rs485"`
		Aiserving bool `json:"aiserving"`
	} `json:"devicemanagement"`
	Monitor struct {
		Telegraf bool `json:"telegraf"`
	} `json:"monitor"`
	Tools struct {
		Nettools struct {
			Ifplugd bool `json:"ifplugd"`
		} `json:"nettools"`
		Mcutools struct {
			Watchdog        bool `json:"watchdog"`
			Powermanagement bool `json:"powermanagement"`
			Mcuserver       bool `json:"mcuserver"`
		} `json:"mcutools"`
		Nodetools struct {
			Cleaner      bool `json:"cleaner"`
			Usblistener  bool `json:"usblistener"`
			Filelistener bool `json:"filelistener"`
		} `json:"nodetools"`
	} `json:"tools"`
	Synctools struct {
		Vpn    bool `json:"vpn"`
		Db     bool `json:"db"`
		Tsdb   bool `json:"tsdb"`
		Mq     bool `json:"mq"`
		Fs     bool `json:"fs"`
		Config bool `json:"config"`
	} `json:"synctools"`
}
