package main

import (
	"encoding/json"
	"io/ioutil"
)

func ProcessJSON(dst string, opts interface{}) error {
	bs, err := json.Marshal(api)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dst, bs, 0600)
	return err
}
