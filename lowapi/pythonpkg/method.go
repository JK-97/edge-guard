package pythonpkg

import (
    "fmt"
    "io/ioutil"
    log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
    "jxcore/lowapi/utils"
    "net/http"
    "os/exec"
    "regexp"
    "strings"
)

func NewPkgClient() PkgClient {
    c := PkgClient{}
    return c
}

//CurPkg CurPkg
func (c *PkgClient) CurPkg() ([]PythonPkg) {
    var err error
    cmdOUt, err := exec.Command("pip3", "freeze").Output()
    if err != nil {
        log.Error(err)
    }
    cmdOUtstr := strings.Replace(string(cmdOUt), " ", "", 1)
    res := strings.Split(cmdOUtstr, "\n")
    temp := make([]PythonPkg, 0)
    for _, per_pkg := range res[:len(res)-1] {
        info := strings.Split(per_pkg, "==")
        temp = append(temp, PythonPkg{Name: info[0], Version: info[1]})

    }
    c.AllPkg = temp

    return c.AllPkg
}

func (c *PkgClient) Internal() ([]PythonPkg, error) {
    allpkg := c.CurPkg()
    allpkgname := make([]string, 0)
    internalpkgname := make([]string, 0)
    temp := make([]PythonPkg, 0)
    for _, pkg := range allpkg {
        allpkgname = append(allpkgname, pkg.Name)
    }

    //get web python pkg
    resp, err := http.Get("http://pypi.jiangxingai.com/simple/")
    if err != nil {
        return nil, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    r := regexp.MustCompile("<a href=\"(.*)/\">")
    web_pkg := r.FindAllStringSubmatch(string(body), -1)
    for _, pkg := range web_pkg {
        internalpkgname = append(internalpkgname, pkg[1])
    }

    //interset

    res := utils.Hash(internalpkgname, allpkgname)
    //interpkg := model.InterfaceToString(res)
    fmt.Println(res)
    for _, internalpkg := range res {
        for _, pkg := range allpkg {
            if pkg.Name == internalpkg {
                temp = append(temp, PythonPkg{internalpkg, pkg.Version})
            }

        }

    }
    c.InternalPkg = temp
    return c.InternalPkg, err
}

func (c *PkgClient) DeletePyPkg() error {
    del_pkg, err := c.Internal()
    if err != nil {
        log.Error(err)
    }

    for _, i := range del_pkg {
        uninstallout, err := exec.Command("pip3", "uninstall", "-y", i.Name).Output()
        if err != nil {
            log.Error(err)
        }
        fmt.Println(string(uninstallout))
    }

    //fmt.Println(pip3out)
    return err
}

func (c *PkgClient) RestorePyPkg() {

    if utils.Exists(PythonPkgRestorePath) {
        restore_pkg, err := ioutil.ReadDir(PythonPkgRestorePath)
        if err != nil {
            log.Error(err)
        }

        for _, pkg := range restore_pkg {
            r := regexp.MustCompile("edgebox-")
            if r.MatchString(pkg.Name()) {
                path := string(PythonPkgRestorePath + pkg.Name())
                _, err := exec.Command("pip3", "install", "-i", "http://pypi.jiangxingai.com/simple/", "--trusted-host", "pypi.jiangxingai.com", path).Output()
                if err != nil {
                    log.Error(err)
                }
            }
        }
        for _, pkg := range restore_pkg {

            path := string(PythonPkgRestorePath + pkg.Name())
            _, err := exec.Command("pip3", "install", "-i", "http://pypi.jiangxingai.com/simple/", "--trusted-host", "pypi.jiangxingai.com", path).Output()
            if err != nil {
                log.Error(err)
            }
            log.Info("has install pythonpkg" + pkg.Name())
        }
    }

}
