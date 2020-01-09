package monitor

import (
	"context"
	"fmt"
	"io/ioutil"
	"jxcore/config/yaml"
	"jxcore/lowapi/logger"
	"jxcore/lowapi/system"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/pkg/errors"
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

			err := tryMount(ei.Path(), mountPath, dst)
			if err != nil {
				return err
			}

		}
	default:
	}
	return nil
}

func InitTF() {
	for src, dst := range yaml.Config.MountCfg {
		if _, err := os.Stat(src); err == nil {
			_, fileName := path.Split(src)
			mountPath := "/media/" + fileName
			err := tryMount(src, mountPath, dst)
			if err != nil {
				logger.Info(err)
			}
		}
	}
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
			return fmt.Errorf("没有插卡")
		}
		_, err := os.Stat(mountPath)
		if err != nil {
			logger.Info("创建文件夹", mountPath)
			_ = os.MkdirAll(mountPath, 0755)
		}
		//没有mount 而有插卡
		logger.Info("备份文件", mountPath)
		err, tmpCopied := tmpCopy(srcPath)
		if err != nil {
			return fmt.Errorf("目录有文件，备份文件失败")
		}
		logger.Info("mount 进行中")
		err = mountTfCard(srcPath, mountPath)
		if err != nil {
			return errors.Wrap(err, "mount 失败")
		}
		if tmpCopied {
			logger.Info("恢复文件")
			err = tmpRestore(srcPath)
			if err != nil {
				return errors.Wrap(err, "恢复文件失败")
			}
		}
		err = link(mountPath, linkPath)
		if err != nil {
			return err
		}
		logger.Infof("Mounted media src: %s, dst: %s, link: %s", srcPath, mountPath, linkPath)
		return nil
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
		err = unLink(linkPath)
		if err != nil {
			return err
		}
		logger.Infof("卡拔出   Unmout ")
		return nil
	}
}

func tmpCopy(srcPath string) (error, bool) {
	_, fileName := path.Split(srcPath)
	dir, err := ioutil.ReadDir(path.Join("/media", fileName))
	if err != nil {
		return err, false
	}
	if len(dir) == 0 {
		logger.Info("无文件备份")
		return nil, false
	}

	logger.Info("正在拷贝文件")
	cmdStr := fmt.Sprintf("cp -r %s %s", path.Join("/media", fileName), path.Join("/tmp", fileName))
	err = system.RunCommand(cmdStr)
	if err != nil {
		return errors.Wrap(err, "拷贝文件失败"), false
	}

	cmdStr = fmt.Sprintf("rm -r %s/*", path.Join("/media", fileName))
	err = system.RunCommand(cmdStr)
	if err != nil {
		return errors.Wrap(err, "拷贝文件失败"), false
	}
	return nil, true
}

func tmpRestore(srcPath string) error {
	_, fileName := path.Split(srcPath)

	_, err := os.Stat(path.Join("/tmp", fileName))
	if err != nil {
		return err
	}
	dir, err := ioutil.ReadDir(path.Join("/tmp", fileName))
	if err != nil {
		return nil
	}
	if len(dir) == 0 {
		return nil
	}
	logger.Info("正在恢复文件")
	cmdStr := fmt.Sprintf("cp -r /tmp/%s/* %s", fileName, path.Join("/media", fileName))
	err = system.RunCommand(cmdStr)
	if err != nil {
		return errors.Wrap(err, "恢复文件失败")
	}
	os.RemoveAll(path.Join("/tmp", fileName))

	return nil
}
