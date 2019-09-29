package config

import (
	"io/ioutil"
	"jxcore/app/schema"

	"gopkg.in/yaml.v2"
)

func LoadYaml(path string) (Yamlsetting schema.YamlSchema, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	yaml.Unmarshal(content, &Yamlsetting)
	return
}
