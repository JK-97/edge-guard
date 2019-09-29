package clean

import (
	"fmt"
	"io"
	"io/ioutil"
	"jxcore/app/model"
	"os"
	"path/filepath"
	"testing"
)

func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

func TestDelFile(t *testing.T) {

	tmpDir, err := prepareTestDirTree("dir/to/walk/skip")
	dir := tmpDir + "/dir/to/walk/skip"
	path := tmpDir + "/dir/to/walk/skip/hello.txt"
	f, err := os.OpenFile(path, os.O_CREATE, 0666)
	io.WriteString(f, "hello")
	defer f.Close()
	if model.Exists(dir) {
		fmt.Println("create dir  succes")
	} else {
		fmt.Println("create dir  fail")
	}
	if model.Exists(path) {
		fmt.Println("create files succes")
	} else {
		fmt.Println("create files fail")
	}
	if err != nil {
		fmt.Printf("unable to create test dir tree: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)
	os.Chdir(tmpDir)

	var path_list = []string{path}
	DelFile(path_list)
	if model.Exists(dir) {
		fmt.Println("dir exist succes")
	}
	if model.Exists(path) {
		fmt.Println("files exist fail")
	}

}
