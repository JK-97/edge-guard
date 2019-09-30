package journal

import (
	"fmt"
	"io/ioutil"
	"jxcore/log"
	"os"
	"path/filepath"
	"time"
)

// SinceToDay 当天的第一秒
func SinceToDay() int64 {
	return Today().Unix()
}

// Today 当天的日期
func Today() time.Time {
	now := time.Now().Local()
	y, m, d := now.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, now.Location())
}

func date(day time.Time) int {
	return day.Year()*10000 + int(day.Month())*100 + day.Day()
}

// Collect 收集日志
func Collect(config map[string]interface{}, day time.Time, folder string, metaFolder string) {
	tr := TimeRange{
		Since: day.Unix(),
		Until: time.Now().Unix(),
	}
	dayInt := day.Year()*10000 + int(day.Month())*100 + day.Day()
	for mode, w := range RegisteredWorkers {
		err := w.InitConfig(&tr, config)
		if err != nil {
			log.Println(err)
			continue
		}
		archive := TarLogArchive{
			TimeRange: tr,
		}
		prefix := fmt.Sprintf("%d-%s", dayInt, mode)
		metaPath := filepath.Join(metaFolder, prefix+".meta.json")
		err = archive.OpenArchive(filepath.Join(folder, prefix+".tar"), metaPath)
		if err != nil {
			log.Println(err)
			continue
		}
		defer archive.CloseArchive()

		archive.Meta.Mode = mode

		if w.CanAppend() {
			archive.SeekForAppend()
		}

		c, wg := w.FetchAsync()

		err = archive.AppendLogs(c, wg)

		if err != nil {
			panic(err)
		}

		if f, err := os.Create(metaPath); err == nil {
			defer f.Close()
			archive.Meta.Serialize(f)
		}
	}
}

// Clean 清理过时的日志
func Clean(arcFolder string, metaFolder string, ttl time.Duration) {
	deadLine := time.Now().Add(-ttl).Unix()

	if infos, err := ioutil.ReadDir(arcFolder); err == nil {
		for _, info := range infos {
			mod := info.ModTime()
			if mod.Unix() < deadLine {
				os.Remove(info.Name())
			}
		}
	}

	if infos, err := ioutil.ReadDir(metaFolder); err == nil {
		for _, info := range infos {
			mod := info.ModTime()
			if mod.Unix() < deadLine {
				os.Remove(info.Name())
			}
		}
	}

	for _, w := range RegisteredWorkers {
		if cw, ok := w.(LoggerCleanWorker); ok {
			cw.Clean(ttl)
		}
	}
}

// RunForever 后台采集日志
// rotate-directory []string
func RunForever(config *map[string]interface{}, span time.Duration, folder string, metaFolder string, ttl time.Duration) {
	log.Info("Begin to collect journal")
	last := Today()

	Collect(*config, last, folder, metaFolder)
	Clean(folder, metaFolder, ttl)
	if span < time.Minute {
		span = 20 * time.Minute
	}
	ticker := time.NewTicker(span)
	d := date(last)
	for {
		select {
		case n := <-ticker.C:
			d2 := date(n)
			log.Info("Collect journal")
			if date(n) != d {
				log.Info("Collect & Clean journal")
				tt := Today().Add(-20 * time.Second)
				Collect(*config, tt, folder, metaFolder)
				Clean(folder, metaFolder, ttl)
			}
			Collect(*config, n, folder, metaFolder)
			last = n
			d = d2
			log.Info("Collect journal Finished")
		}
	}

}
