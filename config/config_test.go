package config

import (
	"bufio"
	"fmt"
	"jxcore/log"
	"os"
	"reflect"
	"testing"
)

func StructInfo(o interface{}, fix string) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		val := v.Field(i).Interface()
		t1 := reflect.TypeOf(val)
		if t1.Kind() != reflect.Struct {
			fmt.Printf(fix+"%s  %v \n", f.Name, val)

		}
		if k := t1.Kind(); k == reflect.Struct {

			StructInfo(val, fix+f.Name+"/")

		}

	}
}

func TestConfig(t *testing.T) {
	f, err := os.OpenFile("/etc/hosts", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error(err)
	}
	var buf []string
	scanner := bufio.NewScanner(f)
	//Reading lines
	for scanner.Scan() {
		line := scanner.Text()
		buf = append(buf, line)
	}
	//Writing from second line on the same file
	for s := 1; s < len(buf); s++ {
		f.WriteString("")
	}
	//Commit changes
	f.Sync()
	f.Close()
}
func TestLoadYaml(t *testing.T) {
	yamlsetting, _ := LoadYaml("/home/marshen/go/src/jxcore/settings.yaml")
	StructInfo(yamlsetting, "")
}
func BenchmarkLoadYaml(b *testing.B) {
	yamlsetting, _ := LoadYaml("/home/marshen/go/src/jxcore/settings.yaml")
	StructInfo(yamlsetting, "")
}
