package v1

import (
	"jxcore/app/model/pythonpackage"
	"net/http"
)

var c = pythonpackage.NewPkgClient()

//PythonPkgInfoGet is
func AllPythonPkgInfoGet(w http.ResponseWriter, r *http.Request) {
	resp := c.CurPkg()
	respondSuccessJSON(resp, w, r)
}

func InternalPythonPkgInfoGet(w http.ResponseWriter, r *http.Request) {
	resp, err := c.Internal()
	if err != nil {
		respondResonJSON(err, w, r, "")
	}
	respondSuccessJSON(resp, w, r)
}

//PythonPkgUinstallDelete is
func PythonPkgUinstallDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		err := c.DeletePyPkg()
		if err != nil {
			log.Error(err)
		}
		respondSuccessJSON(nil, w, r)
	} else {
		respondResonJSON(nil,w,r,"method not support")
	}


}

//PythonPkgURestorePost is
func PythonPkgURestorePost(w http.ResponseWriter, r *http.Request) {
	go c.RestorePyPkg()
	respondSuccessJSON(nil, w, r)
}
