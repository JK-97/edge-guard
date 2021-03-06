package ceph

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/JK-97/edge-guard/lowapi/logger"
)

// CephUmount 取消 ceph 的挂载
func CephUmount() {
	cmd := exec.Command("fusermount", "-u", "-z", cephMountPath)
	cmd.Run()
}

// Cephmount 挂载 ceph
func Cephmount() {
	CephUmount()

	var tmpPath string // 作为保存临时目录，临时保存文件

	infos, err := ioutil.ReadDir(cephMountPath)
	if err != nil && os.IsNotExist(err) {
		// 目录不存在，可以挂载
		os.MkdirAll(cephMountPath, 0666)
	}
	if len(infos) > 0 {
		tmpPath, err := ioutil.TempDir(filepath.Dir(cephMountPath), "remote")
		if err == nil {
			log.Info("Mount Ceph", err)
		}
		cmd := exec.Command("mv", cephMountPath, tmpPath)
		err = cmd.Run()
		if err != nil {
			log.Error("Move out", err)
			tmpPath = ""
		}
		os.Mkdir(cephMountPath, 0666)
	}

	cephcmd := exec.Command("/bin/sh", "-c", "ceph-fuse -m master:6789  /data/edgebox/remote")
	cephcmd.Stdout = os.Stdout
	cephcmd.Stderr = os.Stdout
	err = cephcmd.Run()
	if err == nil && tmpPath != "" {
		// 将临时目录中的文件，移动会 Ceph 目录
		cmd := exec.Command("mv", strings.TrimSuffix(tmpPath, "/")+"/remote", filepath.Dir(cephMountPath))
		err := cmd.Run()
		if err == nil {
			os.Remove(tmpPath)
		} else {
			log.Error("Move back", cmd.Args, err)
		}
	}
}

func TmpFsMount() error {
	_ = os.RemoveAll(tmpfsPath)
	err := os.Mkdir(tmpfsPath, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = syscall.Mount("tmpfs", "/data/tmpfs", "tmpfs", uintptr(syscall.MS_NOSUID|syscall.MS_NOEXEC|syscall.MS_NODEV), "")
	return err
}
func TempUmount() error {
	return syscall.Unmount("/data/tmpfs", 0)
}

func EnsureTmpFs() error {
	log.Info("Checking tmpfs")
	rawdata, err := ioutil.ReadFile(fstabFilePath)
	if err != nil {
		return err
	}
	if !strings.Contains(string(rawdata), fstabRecord) {
		if err := TmpFsMount(); err != nil {
			return err
		}
		output := string(rawdata) + "\n" + fstabRecord
		if err = ioutil.WriteFile(fstabFilePath, []byte(output), 0644); err != nil {
			return err
		}
		log.Info("Mount tmpfs success")
	}
	return nil
}
