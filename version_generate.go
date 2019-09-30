// +build ignore

// generates version/version_impl.go.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const (
	versionTemplate = `
package version

func init() {
	GitCommit = "{{ .commit }}"
	Version   = "{{ .version }}"
	BuildDate = "{{ .date }}"	
}
`
)

func main() {
	var (
		version = "v1.0.0"
		commit  = "?"
	)

	tpl := template.Must(template.New("version").Parse(versionTemplate))

	git, err := exec.LookPath("git")
	if err == nil {
		gitCmd := exec.Command(git, "rev-parse", "--short", "HEAD")
		if out, err := gitCmd.Output(); err == nil {
			commit = strings.TrimSpace(string(out))
		} else {
			fmt.Println(err)
		}

		gitCmd = exec.Command(git, "tag", "--list")
		if out, err := gitCmd.Output(); err == nil {
			s := strings.TrimSpace(string(out))
			tags := strings.Split(s, "\n")
			if len(tags) > 0 {
				version = tags[len(tags)-1]
			}
		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}

	data := map[string]string{
		"version": version,
		"commit":  commit,
		"date":    time.Now().Format(time.RFC3339),
	}

	fp, err := os.Create("version/version_impl.go")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	if err := tpl.Execute(fp, data); err != nil {
		panic(err)
	}
}
