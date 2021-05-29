package goapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var Basepath string = ""
var Host string = ""
var ExtraHeaders map[string]string = make(map[string]string)

func invoke(m string, path string, bodyo interface{}) (*json.Decoder, error) {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(bodyo)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b.Bytes())
	req, err := http.NewRequest(m, Host+Basepath+path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")

	for k, v := range ExtraHeaders {
		req.Header.Set(k, v)
	}

	cli := http.Client{}
	res, err := cli.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		return nil, errors.New(string(bs))
	}

	ret := json.NewDecoder(res.Body)
	return ret, nil
}

func SomeAPI(req []*time.Time) (res string, err error) {
	var dec *json.Decoder
	dec, err = invoke("POST", "/someapi", req)
	if err != nil {
		return
	}
	var ret string
	err = dec.Decode(&ret)
	return ret, err
}
