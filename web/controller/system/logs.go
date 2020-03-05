package system

import (
	"fmt"
	"jxcore/lowapi/logger"
	"jxcore/oplog"
	"jxcore/oplog/types"
	"jxcore/web/controller/utils"
	"net/http"
	"strconv"
	"time"
)

func GetOplog(w http.ResponseWriter, r *http.Request) {
	type responce struct {
		Oplogs []types.Oplog `json:"oplogs"`
		Total  int           `json:"total"`
	}
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
	if err != nil {
		utils.RespondReasonJSON(nil, w, err.Error(), 400)
		return
	}
	offset, limit := utils.GetPageInfo(r)
	length := len(findResult)

	if offset > length {
		findResult = []types.Oplog{}
	} else if offset+limit < length {
		findResult = findResult[offset : offset+limit]
	} else if offset+limit > length {
		findResult = findResult[offset:]
	}

	resp := &responce{
		Oplogs: findResult,
		Total:  length,
	}

	utils.RespondSuccessJSON(resp, w)
}

func DownloadOplog(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, oplog.GetLogFileName())
}

func ListFilter(offset, limit int, data []types.Oplog) ([]types.Oplog, int, error) {

	res := make([]types.Oplog, 0)
	length := len(res)
	if offset > length {
		return []types.Oplog{}, length, fmt.Errorf("offset more than total")
	}

	if offset+limit < length {
		return res[offset : offset+limit], length, nil

	}
	return res[offset:], length, nil
}
