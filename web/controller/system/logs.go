package system

import (
	"jxcore/lowapi/logger"
	"jxcore/oplog"
	"jxcore/web/controller/utils"
	"net/http"
	"strconv"
	"time"
)

func GetOplog(w http.ResponseWriter, r *http.Request) {
	querys := r.URL.Query()
	fromStr, ok := querys["from"]
	if !ok {
		utils.RespondReasonJSON(nil, w, "notfound args from", 400)
		return
	}
	untilStr, ok := querys["until"]
	if !ok {
		utils.RespondReasonJSON(nil, w, "notfound args util", 400)
		return
	}
	logMessageType, ok := querys["type"]
	if !ok {
		utils.RespondReasonJSON(nil, w, "notfound args util", 400)
		return
	}
	from, err := strconv.ParseInt(fromStr[0], 10, 64)
	if err != nil {
		utils.RespondReasonJSON(nil, w, "invalid args from", 400)
		return
	}
	until, err := strconv.ParseInt(untilStr[0], 10, 64)
	if err != nil {
		utils.RespondReasonJSON(nil, w, "invalid args util", 400)
		return
	}
	untilTime := time.Unix(until, 0)
	fromTime := time.Unix(from, 0)
	logger.Info(untilTime.Format("2006-01-02 15:04:05"), fromTime.Format("2006-01-02 15:04:05"))

	findResult, err := oplog.FindMany(oplog.DefaultTimeFilter(fromTime, untilTime), oplog.DefaultTypeFilter(logMessageType[0]))
	// logger.Info(findResult)
	if err != nil {
		utils.RespondReasonJSON(nil, w, err.Error(), 400)
		return
	}
	utils.RespondSuccessJSON(findResult, w)
}

func DownloadOplog(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, oplog.GetLogFileName())
}
