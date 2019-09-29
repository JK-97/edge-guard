package utils

import (
	"archive/zip"
	"bytes"
	"io"
	"jxcore/log"
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
