package clean

import (
	"jxcore/log"
	"jxcore/utils"
	"os"
	"path/filepath"
)

//DelFile is
func DelFile(path_list []string) {
	//Clean up all files under the directory, but save the folder structure ,
	for _, per_path := range path_list {
		if utils.Exists(per_path) {
			if utils.IsDir(per_path) {
				filepath.Walk(per_path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						log.Infof("prevent panic by handling failure accessing a path %q: %v", path, err)
						return err
					}
					if !info.IsDir() {
						os.Remove(path)
						log.Info("remove path : ", path)
						return nil
					}
					return nil
				})
			} else {
				os.Remove(per_path)
				log.Info("remove path : ", per_path)
			}
		}
	}
}

//ResetFile is
func ResetFile(path_list []string) {
	//Clean up all files under the directory, but save the files structure ,
	for _, per_path := range path_list {
		if utils.Exists(per_path) {
			if utils.IsDir(per_path) {
				filepath.Walk(per_path, func(path string, info os.FileInfo, err error) error {

					if err != nil {
						log.Infof("prevent panic by handling failure accessing a path %q: %v", path, err)
						return err
					}
					if !info.IsDir() {
						f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
						defer f.Close()
						if err != nil {
							log.Info(err)
						}
						f.WriteString("")
						log.Info("reset file : ", path)
						return nil
					}
					return nil
				})
			} else {
				f, err := os.OpenFile(per_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				defer f.Close()
				if err != nil {
					log.Error(err)
				}
				f.WriteString("")
				log.Info("reset path: : ", per_path)
			}
		}
	}
}
