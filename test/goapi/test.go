package goapi

import (
	"context"
	"crypto"
	"log"
	"time"
)

/*@API*/
type AStr struct {
	Country       string
	City          string
	HouseNumber   int64
	IsCondo       bool
	SomeWeirdTest string `json:"SUPERCALIFRAGILISPEALIDOUX"`
	Recursive     map[string]AStr
	Arrofpstr     []string `json:"arrofpstr,omitempty"`
	When          time.Time
	Some          crypto.Decrypter
}

/*
@API
@PATH: /someapi
@PERM: ASD
@VERB: POST
*/
func SomeAPI(ctx context.Context, s []time.Time) (out string, err error) {
	log.Printf("Got: %#v", s)
	out = time.Now().String() + " - Hey Ya!"
	return
}

/*
@API
@PATH: /someapi2
@PERM: ASD
@VERB: GET
*/
func SomeGET(ctx context.Context, s []*AStr) (out string, err error) {
	print("Got:" + s[0].SomeWeirdTest)
	out = time.Now().String() + " - Hey Ya!"
	return
}

//
///*
//@API
//@PATH: /someapi
//@PERM: ASD
//@VERB: PUT
//*/
//func SomePUT(ctx context.Context, s string) (out string, err error) {
//	print("Got:" + s)
//	out = time.Now().String() + " - Hey Ya!"
//	return
//}
//
///*
//@API
//@PATH: /someapi
//@PERM: ASD
//@VERB: DELETE
//*/
//func SomeAPI2(ctx context.Context, s *crypto.Hash) ([]string, error) {
//	return nil, nil
//}
