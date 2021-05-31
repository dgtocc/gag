package main

import (
	cliapi "github.com/dgtocc/gag/test/gocli/api"
	"log"
)

func main() {
	cliapi.Host = "http://localhost:8080"
	res, err := cliapi.ApiMethod01(&cliapi.ASimpleReq{Data: "ASD"})
	if err != nil {
		panic(err)
	}
	log.Printf("%#v", res)
}
