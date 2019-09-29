package v1

import (
	"jxcore/app/model/mongo"
	"net/http"
)

func MongoinfoGet(w http.ResponseWriter, r *http.Request) {

}

func MongoRemoveDelete(w http.ResponseWriter, r *http.Request) {
	mongo.UnInstallMongo()
	respondSuccessJSON(nil, w, r)
}

func MongoRestorePost(w http.ResponseWriter, r *http.Request) {
	mongo.InstallMongo()
	respondSuccessJSON(nil,w,r)
}
