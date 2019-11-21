package yaml

import (
	"io/ioutil"
	"jxcore/management/programmanage"
	"reflect"

	"gopkg.in/yaml.v2"
)

func init() {
	err := loadYaml(configPath)
	if err != nil {
		panic(err)
	}
}

func loadYaml(path string) (err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, Config)
	return
}

// ParseAndCheck 通过递归调用获取子类型的信息
func ParseAndCheck(o interface{}, fix string) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		t1 := reflect.TypeOf(val)
		if t1.Kind() == reflect.String {
			continue
		}

		if t1.Kind() != reflect.Struct {
			//path := strings.ToLower(fix + f.Name)
			if b, ok := val.(bool); ok {
				if b {
					//binfile := strings.ToLower("/edge/" + path + "/bin/" + f.Name)
					//if strings.Count(binfile, "synctools") != 0 {
					//    binfile = strings.ReplaceAll(binfile, "synctools", "mnt")
					//}
					programmanage.AddDependStart(f.Name)

				} else {
				}
			}
		}

		if k := t1.Kind(); k == reflect.Struct {
			ParseAndCheck(val, fix+f.Name+"/")

		}
	}

}
