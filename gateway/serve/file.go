package serve

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileHandler 处理文件
type FileHandler struct {
}

// checkPath 判断目录是否存在
func checkPath(source string, destination *string) error {
	var err error
	if _, err = os.Stat(source); err != nil {
		return err
	}

	if strings.HasSuffix(*destination, "/") {
		if _, err = os.Stat(*destination); err != nil {
			if os.IsNotExist(err) {
				if err = os.MkdirAll(*destination, 0666); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		*destination = filepath.Join(*destination, filepath.Base(source))
		return nil
	}

	dir := filepath.Dir(*destination)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0666)
		}
		return err
	}
	return nil
}

func (h *FileHandler) copyFile(source string, destination string) error {
	ifile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer ifile.Close()
	ofile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer ofile.Close()
	_, err = io.Copy(ofile, ifile)

	return err
}

// SyncFile 同步文件
func (h *FileHandler) SyncFile(source string, destination string) error {
	if err := checkPath(source, &destination); err != nil {
		return err
	}

	return h.copyFile(source, destination)
}

// SyncFolder 同步目录
func (h *FileHandler) SyncFolder(source string, destination string) error {
	if err := checkPath(source, &destination); err != nil {
		return err
	}
	cmd := exec.Command("cp", "-ur", source, destination)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + "\n" + string(b))
	}

	return nil
}

type fileSyncRequest struct {
	Source      string `json:"src,omitempty"`
	Destination string `json:"dst,omitempty"`
	IsDir       bool   `json:"dir,omitempty"`
}

// SyncFileSystemHTTP 文件同步
func (h *FileHandler) SyncFileSystemHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req fileSyncRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Source == "" || req.Destination == "" {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.IsDir {
		err = h.SyncFolder(req.Source, req.Destination)
	} else {
		err = h.SyncFile(req.Source, req.Destination)
	}
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteSucess(w)
}

type cleanJournalRequest struct {
	Size string `json:"size"`
	Time string `json:"time"`
}

// CleanJournalHTTP 清理日志
func (h *FileHandler) CleanJournalHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req cleanJournalRequest
	err = json.Unmarshal(b, &req)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	args := make([]string, 0, 2)

	if req.Size != "" {
		args = append(args, "journalctl --vacuum-size="+req.Size)
	}
	if req.Time != "" {
		args = append(args, "journalctl --vacuum-time="+req.Time)
	}

	if len(args) == 0 {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c := strings.Join(args, " & ")

	cmd := exec.Command("sh", "-c", c)

	buff, err := cmd.CombinedOutput()
	if err != nil {
		Error(w, string(buff)+err.Error(), http.StatusInternalServerError)
		return
	}

	outPut := map[string]interface{}{
		"output": string(buff),
	}
	WriteData(w, &outPut)
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorWithCode(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/clean/journal" {
		h.CleanJournalHTTP(w, r)
		return
	}

	h.SyncFileSystemHTTP(w, r)
}
