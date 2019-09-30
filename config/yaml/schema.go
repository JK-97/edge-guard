package yaml

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)




func LoadYaml(path string) (Yamlsetting schema.YamlSchema, err error) {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return
    }
    yaml.Unmarshal(content, &Yamlsetting)
    return
}
