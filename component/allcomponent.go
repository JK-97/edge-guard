package component

import (
	"hash/fnv"
	"jxcore/component/process"
	"jxcore/config"
	"jxcore/log"
	"jxcore/utils"
	"reflect"
	"strings"
	"sync"
)

type ComponentPID struct {
	Gpid []*process.Process
}

var ComponentPidInfo ComponentPID

var ComponentPath = map[string]string{}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

//通过递归调用获取子类型的信息
func PraseComponentCfg(o interface{}, fix string) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		t1 := reflect.TypeOf(val)
		if t1.Kind() != reflect.Struct {
			path := strings.ToLower(fix + f.Name)
			ComponentPath[strings.ToLower(f.Name)] = strings.ToLower("/edge/" + path)
			if b, ok := val.(bool); ok {
				if b {
					binfile := strings.ToLower("/edge/" + path + "/bin/" + f.Name)
					if strings.Count(binfile, "synctools") != 0 {
						binfile = strings.ReplaceAll(binfile, "synctools", "mnt")
					}
					if utils.Exists(binfile) {
						// start the component
						//cmd := exec.Command(binfile, "-c", binfile+".cfg")
						//cmd.Start()

						p := process.NewProcess(string(hash(binfile)), binfile)
						p.Start(true)
						ComponentPidInfo.Gpid = append(ComponentPidInfo.Gpid, p)
						//log.WithFields(log.Fields{"component": f.Name}).Info("up "+"at pid : ", p.GetCmd().Process.Pid)

					} else {
						log.WithFields(log.Fields{"component": f.Name}).Warn("not found path ")
					}

				} else {
				}
			}
		}

		if k := t1.Kind(); k == reflect.Struct {
			PraseComponentCfg(val, fix+f.Name+"/")

		}
	}

}

func ComponentEmiter() {
	cfg := config.Config()
	yamlpath := cfg.GetString("yamlsettings")
	yamlsetting, err := config.LoadYaml(yamlpath)
	if err != nil {
		log.Error(err)
	}
	PraseComponentCfg(yamlsetting, "")

}

// StopComponent 停止组件
func StopComponent(component string) {
	wg := new(sync.WaitGroup)
	for _, perprocess := range ComponentPidInfo.Gpid {
		if component == "" || perprocess.GetName() == component {
			wg.Add(1)
			go func(p *process.Process) {
				defer wg.Done()
				p.Stop(true)
			}(perprocess)
		}
	}
	wg.Wait()
}

// StartComponent 运行组件
func StartComponent(component string) {
	wg := new(sync.WaitGroup)
	for _, perprocess := range ComponentPidInfo.Gpid {
		if component == "" || perprocess.GetName() == component {
			wg.Add(1)
			go func(p *process.Process) {
				defer wg.Done()
				p.Start(true)
			}(perprocess)
		}
	}
	wg.Wait()
}

//
//func CleanerAdministrator(){
//	sigc := make(chan os.Signal, 1)
//	signal.Notify(sigc,
//		syscall.SIGHUP,
//		syscall.SIGINT,
//		syscall.SIGTERM,
//		syscall.SIGQUIT)
//	go func() {
//		_ = <-sigc
//		for  _,perone :=range ComponentPidInfo.Gpid{
//			perone.Stop(false)
//		}
//		mainprocess,err:=os.FindProcess(os.Getpid())
//		if err != nil {
//			log.Error(err)
//		}
//		mainprocess.Kill()
//
//
//}()

//
//}
