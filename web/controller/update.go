package controller

import (
    "encoding/json"
    "io/ioutil"
    "jxcore/core"
    "jxcore/core/device"
    "jxcore/log"
    "jxcore/lowapi/utils"
    "net/http"
    "jxcore/management/updatemanage"
)

func UpdateByDeb(w http.ResponseWriter, r *http.Request) {
    updateprocess := updatemanage.GetUpdateProcess()
    log.Warn("hello")
    if updateprocess.GetStatus() != updatemanage.FINISHED {
        w.WriteHeader(400)
        respondResonJSON(nil, w, r, "machine is busy to updating,please update later")
    } else {
        reqrawdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Error(err)
        }
        reqinfo := updatemanage.Reqdatastruct{}
        err = json.Unmarshal(reqrawdata, &reqinfo)
        if err != nil {
            log.Error(err)
        }
        indentdata, err := json.MarshalIndent(reqinfo.Data, "", "  ")

        utils.CheckErr(err)

        updateprocess.SetNewTarget(indentdata)

        go func() {
            core.UpdateCore(30)
        }()
        deviceinfo, err := device.GetDevice()
        if err != nil {
            log.Error(err)
        }
        respdata := updatemanage.Respdatastruct{
            Status:   updatemanage.FINISHED.String(),
            WorkerId: deviceinfo.WorkerID,
            PkgInfo:  updatemanage.ParseVersionFile(),
        }
        respondJSON(respdata, w, r)
        //wg.Add(2)

    }

}
