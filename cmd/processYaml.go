package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ProcessYaml(dst string, opts interface{}) error {
	bs, err := yaml.Marshal(api)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, bs, 0600)
	return err
}
