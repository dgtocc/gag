package gaglib

import (
	"log"
)

func Err(e error) {
	if e != nil {
		panic(e)
	}
}

func Debug(s string, p ...interface{}) {
	log.Printf("DEBUG: "+s, p...)
}

func Log(s string, p ...interface{}) {
	log.Printf("LOG: "+s, p...)
}
