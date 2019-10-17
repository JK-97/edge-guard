// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
    "jxcore/log"
    "jxcore/lowapi/network"
    "jxcore/lowapi/utils"
    "net/http"
    "strings"

    // 调试
    _ "net/http/pprof"
    "os/exec"

    "github.com/spf13/cobra"
)

// serveCmd represents the serve command
var checkCmd = &cobra.Command{
    Use:   "serve",
    Short: "Serve http backend for jxcore",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        if _, err := exec.LookPath("docker"); err == nil {
            log.Info("INSTALLED", "docker")
        }

        if _, err := exec.LookPath("supervisorD"); err == nil {
            log.Info("INSTALLED", "supervisorD")
        }
        if _, err := exec.LookPath("openvpn "); err == nil {
            log.Info("INSTALLED", "openvpn")
        }
        if _, err := exec.LookPath("docker-compose"); err == nil {
            log.Info("INSTALLED", "docker-compose")
        }
        if _, err := exec.LookPath("aptitude"); err == nil {
            log.Info("INSTALLED", "aptitude")
        }
        if _, err := exec.LookPath("dnsmasq"); err == nil {
            log.Info("INSTALLED", "dnsmasq")
        }

        out, err := exec.Command("/bin/bash", "-c", "pip3 list | grep edgebox").Output()
        utils.CheckErr(err)
        if strings.Contains(string(out), "edgebox") {
            log.Info("INSTALLED", "edgebox sdk")
        }

        //dockerInstance:=docker.NewClient()
        //dockerList,err:=dockerInstance.ImagesList()
        //rawdata,err:=ioutil.ReadFile(restoreImagePath+"/desc.json")
        //utils.CheckErr(err)
        //var dockerinfo =make(map[string]map[string]string)
        //json.Unmarshal(rawdata,&dockerinfo)
        //currentImage := make(map[string]string,0)
        //for _,dockerimage :=range dockerList{
        //    currentImage[dockerimage.ImageID]=dockerimage.Tag[0]
        //}
        //
        //for _,imageinfo := range dockerinfo{
        //    if currentImage[imageinfo["id"]] !=""{
        //        log.Info("INSTALLED IMAGE",imageinfo["repo"])
        //    }else {
        //        log.Info("NOT INSTALLED IMAGE",imageinfo["repo"])
        //    }
        //}
        network.CheckNetwork()

        if out, _ := exec.Command("pgrep", "jxcore_service").Output(); len(out) != 0 {
            log.Info("has detect the jxcore running")
            resp, err := http.Get("http://localhost:80/ping")
            utils.CheckErr(err)
            if resp.StatusCode == 200 {
                log.Info("jxcore web service health")
            }
        }
        for programname , shellpath :=range BinFilesMAP {
            if out, _ := exec.Command("pgrep", "-f", shellpath).Output(); len(out) != 0 {
                log.Info(programname+ " running ")
            }else {
                log.Info(programname+ " not running ")
            }
        }
    },
}

func init() {
    rootCmd.AddCommand(checkCmd)

}
