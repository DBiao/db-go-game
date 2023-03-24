package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

func YamlToStruct(file string, out interface{}) (err error) {
	var buf []byte
	buf, err = os.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(buf, out)
	return
}
