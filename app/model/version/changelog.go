package version

import (
	"encoding/json"
	"io/ioutil"
	"jxcore/config"
	"jxcore/log"
	"os"
	"time"
)

func ChangLog(praseres map[string]string) map[string]map[string]map[string]string {
	datestring := time.Now().Format("2006-01-01 15:04:05")
	data, _ := ioutil.ReadFile(config.InterSettings.ChangeLog)
	olderinfo := make(map[string]map[string]map[string]string)
	json.Unmarshal(data, &olderinfo)
	var date string = ""
	for olderdate, conponentinfo := range olderinfo {
		for edgeversion, _ := range conponentinfo["edge"] {

			if edgeversion == praseres["edge"] {
				date = olderdate
			}
		}

	}

	if date != "" {
		theversion := olderinfo[date]
		for name, ver := range praseres {
			if theversion[name] == nil {
				theversion[name] = map[string]string{ver: datestring}
			} else {
				if theversion[name][ver] == "" {
					theversion[name][ver] = datestring
				}
			}
		}

		olderinfo[date] = theversion
	} else {
		theversion := make(map[string]map[string]string)
		for name, version := range praseres {
			theversion[name] = map[string]string{version: datestring}
		}
		olderinfo[datestring] = theversion
	}

	jsonfile, _ := json.MarshalIndent(olderinfo, "", "    ")
	f, err := os.OpenFile("changelog.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer f.Close()
	if err != nil {
		log.Error(err)
	}
	f.WriteString(string(jsonfile))
	return olderinfo

}
