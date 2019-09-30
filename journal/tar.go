package journal

import (
	"archive/tar"
	"io"
	"os"
	"sync"
	"time"
)

const seekOffset int64 = 1 << 10
const dateFormat string = "20060102"

// TarLogArchive tar 格式的文件包
type TarLogArchive struct {
	TimeRange
	KeepEmptyLog bool // 保留大小为 0 的日志
	Meta         LogArchiveMeta
	writer       *tar.Writer
	f            *os.File
}

// OpenArchive 打开文件包
func (a *TarLogArchive) OpenArchive(archive, meta string) error {
	pMeta := &a.Meta
	if f, err := os.Open(meta); err == nil {
		defer f.Close()
		a.Meta.Deserialize(f)
		a.Since = a.Meta.ModifyTime
	} else {
		a.Since = SinceToDay()
		y, m, d := Today().Date()
		pMeta.Day = y*10000 + int(m)*100 + d
		pMeta.CreateTime = a.Since
	}
	if a.Until == 0 {
		a.Until = time.Now().Unix()
	}

	if pMeta.Rotates == nil {
		pMeta.Rotates = make(map[string]MetaFileItem)
	}

	f, err := os.OpenFile(archive, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	a.f = f

	// stat, err := f.Stat()
	// if stat.Size() >= seekOffset {
	// 	if _, err = f.Seek(-seekOffset, io.SeekEnd); err != nil {
	// 		return err
	// 	}
	// }
	a.writer = tar.NewWriter(f)

	return nil
}

// SeekForAppend 为追加写入做准备
func (a *TarLogArchive) SeekForAppend() error {
	f := a.f
	stat, err := a.f.Stat()
	if stat.Size() >= seekOffset {
		if _, err = f.Seek(-seekOffset, io.SeekEnd); err != nil {
			return err
		}
	}
	return nil
}

// CloseArchive 关闭压缩包
func (a *TarLogArchive) CloseArchive() error {
	return a.writer.Close()
}

// AppendLogs 添加日志
func (a *TarLogArchive) AppendLogs(c chan BufferedLoggerWrapper, wg *sync.WaitGroup) error {
	if a.Until == 0 {
		a.Until = time.Now().Unix()
	}
	cTime := time.Now().Local()
	go func() {
		for w := range c {
			if w == nil {
				wg.Done()
				continue
			}
			b := w.Bytes()
			size := int64(len(b))
			if size == 0 && (!a.KeepEmptyLog) {
				wg.Done()
				continue
			}
			name := w.Name()

			hdr := tar.Header{
				Name:       name,
				Size:       size,
				Mode:       0666,
				AccessTime: cTime,
				ChangeTime: cTime,
				Format:     tar.FormatGNU,
			}
			a.writer.WriteHeader(&hdr)
			a.writer.Write(b)
			a.Meta.Files = append(a.Meta.Files, name)
			a.Meta.Rotates[name] = MetaFileItem{
				a.TimeRange,
			}
			wg.Done()
		}
	}()

	wg.Wait()
	close(c)
	a.Meta.ModifyTime = a.Until
	return a.writer.Flush()
}
