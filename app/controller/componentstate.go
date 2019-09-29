package controller

import (
	"jxcore/component"
	"net/http"
)

func ComponentState(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("component")

	if name != "" {
		var data = map[string]string{}
		for _, perprocess := range component.ComponentPidInfo.Gpid {
			if name == perprocess.GetName() {
				//reader := bufio.NewReader(perprocess.StdoutLog)
				//buf := make([]byte, 1024)
				//
				//for {
				//
				//	line, err := reader.Read(buf)
				//	if err != nil || io.EOF == err {
				//		break
				//	}
				//
				//	fmt.Println(buf[:line])
				//	if line !=1024{
				//		break
				//	}
				//}
				data[perprocess.GetName()] = perprocess.GetState().String()
			}

		}
		respondSuccessJSON(data, w, r)
	} else {
		var data = map[string]string{}
		for _, perprocess := range component.ComponentPidInfo.Gpid {

			data[perprocess.GetName()] = perprocess.GetState().String()

		}

		respondSuccessJSON(data, w, r)
	}

}
