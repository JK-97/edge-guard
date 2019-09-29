package version

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"jxcore/log"
	"os"
	"path/filepath"
)

type ComponentInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func PraseVersionFile() (map[string]string) {
	versioninfo := map[string]string{}
	err:=filepath.Walk("/edge", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			if len(path)>7 {

				if (path[len(path)-7:] == "version") {

					filebyte, err := ioutil.ReadFile(path)
					if err != nil {
						log.Error(err)
					}
					componentinfo := ComponentInfo{}
					yaml.Unmarshal(filebyte, &componentinfo)
					versioninfo[componentinfo.Name] = componentinfo.Version
				}
			}


		}
		return err
	})
	if err !=nil{
		log.Error(err)
	}
	return versioninfo
}
