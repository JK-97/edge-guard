package utils

import (
	"archive/zip"
	"bytes"
	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"io"
	"os"
	"path/filepath"
)

func Unzip(bytefile []byte, target string) error {
	a := bytes.NewReader(bytefile)
	reader, err := zip.NewReader(a, int64(len(bytefile)))
	if err != nil {
		log.Error(err)
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		log.Error(err)
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			log.Error(err)
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			log.Error(err)
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			log.Error(err)
		}
	}
	return nil
}

func SaveFile(tempfilename string, binfile io.Reader) error {
	fW, err := os.Create(tempfilename)
	if err != nil {
		log.Error(tempfilename + "文件创建失败")
		return err
	}
	defer fW.Close()
	_, err = io.Copy(fW, binfile)
	if err != nil {
		log.Error(tempfilename + "文件保存失败")
		return err
	}
	return err
}

// Exists Exists
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//DelFile is
func DelFile(path_list []string) {
	//Clean up all files under the directory, but save the folder structure ,
	for _, per_path := range path_list {
		if Exists(per_path) {
			if IsDir(per_path) {
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
		if Exists(per_path) {
			if IsDir(per_path) {
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
