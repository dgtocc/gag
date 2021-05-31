package goapi

import (
	"context"
	"crypto"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

/*@API*/
type ASimpleReq struct {
	Data string
}

/*@API*/
type ASimpleRes struct {
	Data string
}

/*@API*/
type AComplexReq struct {
	Country       string
	City          string
	HouseNumber   int64
	IsCondo       bool
	SomeWeirdTest string `json:"SUPERCALIFRAGILISPEALIDOUX"`
	Recursive     map[string]AComplexReq
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
func ApiMethod01(ctx context.Context, s *ASimpleReq) (out *ASimpleRes, err error) {
	log.Printf("Got: %#v", s)
	out = &ASimpleRes{}
	out.Data = time.Now().String() + " - Hey Ya!"
	return
}

/*
@API
@PATH: /someapi2
@PERM: ASD
@VERB: GET
*/
func ApiMethod02(ctx context.Context, s *AComplexReq) (out *ASimpleRes, err error) {
	gctx := ctx.Value("CTX").(*gin.Context)
	log.Printf(gctx.FullPath())
	print("Got:" + s.SomeWeirdTest)
	out = &ASimpleRes{}
	out.Data = time.Now().String() + " - Hey Ya!"
	return
}
