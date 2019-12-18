package monitor

import (
	"context"
	"errors"
	"io/ioutil"
	"jxcore/config/yaml"
	"jxcore/lowapi/logger"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/rjeczalik/notify"
)

/*
1 jxcore起来带起filelistener
2 主动，或被动检测到tf卡插入，并且确保docker运行之前
3 进行mount，
4 检测到 tf拔出，则umount
*/

// 获取挂载配置: map[string]map[string]string 源文件夹->源文件名->目标路径
func GetMountCfg() map[string]map[string]string {
	var mergedCfg = map[string]map[string]string{}
	for src, dst := range yaml.Config.MountCfg {
		dir, filename := path.Split(src)
		if _, ok := mergedCfg[dir]; !ok {
			mergedCfg[dir] = map[string]string{filename: dst}
		} else {
			mergedCfg[dir][filename] = dst
		}
	}
	return mergedCfg
}

// 监听dirPath设备文件夹变化，挂载到目标文件夹
func MountListener(ctx context.Context, dirPath string, mapSrcDst map[string]string) error {
	eventChan := make(chan notify.EventInfo, 2)

	if err := notify.Watch(dirPath, eventChan, notify.Create, notify.Remove); err != nil {
		return err
	}
	defer notify.Stop(eventChan)

	for {
		select {
		case event := <-eventChan:
			if err := fileHandler(event, mapSrcDst); err != nil {
				logger.Error(err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func fileHandler(ei notify.EventInfo, mapSrcDst map[string]string) error {
	switch ei.Event() {
	case notify.Create, notify.Remove:
		_, fileName := path.Split(ei.Path())
		if dst, ok := mapSrcDst[fileName]; ok {
			mountPath := "/media/" + fileName
			err := os.MkdirAll(mountPath, 0755)
			if err != nil {
				return err
			}
			err = tryMount(ei.Path(), mountPath, dst)
			if err != nil {
				return err
			}
			logger.Infof("Mounted media src: %s, dst: %s, link: %s", ei.Path(), mountPath, dst)
		}
	default:
	}
	return nil
}

// mountTfCard 进行mount
func mountTfCard(mountPoint, mountPath string) error {
	return exec.Command("mount", mountPoint, mountPath).Run()
}

// umount
func unmount(mountPath string) error {
	return syscall.Unmount(mountPath, 0)
}

// checkMount 检查 mountPath的mount 状态
// /dev/root / ext4 rw,relatime,data=ordered 0 0
// devtmpfs /dev devtmpfs rw,relatime,size=1966428k,nr_inodes=491607,mode=755 0 0
func checkMount(mountPath string) (bool, error) {
	data, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return false, err
	}
	allData := strings.TrimSpace(string(data))
	lines := strings.Split(allData, "\n")
	for _, line := range lines {
		info := strings.Split(line, " ")
		// 第二个字段为挂在的target
		if info[1] == mountPath {
			return true, nil
		}
	}
	return false, nil
}

//CheckTfCard 检查tf卡路径是否存在
func CheckTfCard(src string) bool {
	_, err := os.Stat(src)
	return err == nil
}

// 创建软链
func link(src, dst string) error {
	_ = os.RemoveAll(dst)
	return exec.Command("ln", "-s", src, dst).Run()
}

// 解除软链
func unLink(dst string) error {
	_ = os.RemoveAll(dst)
	return exec.Command("unlink", dst).Run()
}

func tryMount(srcPath, mountPath, linkPath string) error {
	tfCardOk := CheckTfCard(srcPath)
	mountOk, mountErr := checkMount(mountPath)
	if mountErr != nil {
		return mountErr
	}

	if !mountOk {
		if !tfCardOk {
			// 没有mount但没有插卡
			return errors.New("没有插卡")
		}
		//没有mount 而有插卡
		err := mountTfCard(srcPath, mountPath)
		if err != nil {
			return err
		}
		return link(mountPath, linkPath)
	}

	if tfCardOk {
		//插卡并且已经mount
		return nil
	} else {
		//没插卡但已经mount
		err := unmount(mountPath)
		if err != nil {
			return err
		}
		return unLink(linkPath)
	}
}
