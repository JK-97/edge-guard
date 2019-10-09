package controller

import (
    "encoding/json"
    "io/ioutil"
    "jxcore/core/device"
    "jxcore/log"
    "jxcore/lowapi/utils"
    updatemange "jxcore/management/updatemanage"
    "net/http"
)

func UpdateByDeb(w http.ResponseWriter, r *http.Request) {
    updateprocess := updatemange.GetUpdateProcess()
    if updateprocess.GetStatus() != updatemange.FINISHED {
        w.WriteHeader(400)
        respondResonJSON(nil, w, r, "machine is busy to updating,please update later")
    } else {
     
            log.Error(err)
        }
        indentdata, err := json.MarshalIndent(reqinfo.Data, "", "  ")

        utils.CheckErr(err)

        updateprocess.SetNewTarget(indentdata)
        
        go func() {
            updateprocess.UpdateSource()
            updateprocess.UploadVersion()
        }()
        deviceinfo, err := device.GetDevice()
        if err != nil {
            log.Error(err)
        }
        respdata := updatemange.Respdatastruct{
            Status:   updatemange.FINISHED.String(),
            WorkerId: deviceinfo.WorkerID,
            PkgInfo:  updatemange.ParseVersionFile(),
        }
        respondJSON(respdata, w, r)
        //wg.Add(2)

    }

}
