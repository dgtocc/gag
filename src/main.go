package main

import (
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"log"
)

var knownMethods map[string]bool = make(map[string]bool)
var httpMapper map[string]map[string]string = make(map[string]map[string]string)
var packageName string = "main"

var CLI struct {
	Yaml struct {
		Src   string `arg help:"Source Dir"`
		Fname string `arg help:"File to be generated"`
	} `cmd help:"Gens YAML metamodel"`
	Json struct {
		Src   string `arg help:"Source Dir"`
		Fname string `arg help:"File to be generated"`
	} `cmd help:"Gens Json metamodel"`
	Gin struct {
		Src string `arg help:"Source Dir"`
	} `cmd help:"Gens Gin Server impl"`
	Gocli struct {
		Src string `arg help:"Source Dir"`
		Dst string `arg help:"Dst file"`
	} `cmd help:"Gens Go Cli impl"`
	Pycli struct {
		Src string `arg help:"Source Dir"`
		Dst string `arg help:"Dst file"`
	} `cmd help:"Gens Python Cli impl"`
	Ts struct {
		Src string `arg help:"Source Dir"`
		Dst string `arg help:"Dst file"`
	} `cmd help:"Gens Typescript Cli impl"`
	Http struct {
		Src string `arg help:"Source Dir"`
		Dst string `arg help:"Dst file"`
	} `cmd help:"Gens Http call impl"`
}

func main() {

	var processor func() error
	kong.ConfigureHelp(kong.HelpOptions{
		NoAppSummary: false,
		Summary:      true,
		Compact:      true,
		Tree:         true,
		Indenter:     nil,
	})
	ctx := kong.Parse(&CLI)
	var err error
	var src string
	switch ctx.Command() {
	case "yaml <src> <fname>":
		log.Printf("Gens YAML")
		src = CLI.Yaml.Src
		processor = func() error {
			return processYaml(CLI.Yaml.Fname, nil)
		}
	case "json <src> <fname>":
		log.Printf("Gens JSON")
		src = CLI.Json.Src
		processor = func() error {
			return processJSON(CLI.Json.Fname, nil)
		}

	case "gin <src>":
		log.Printf("Gen Gin Server")
		src = CLI.Gin.Src
		processor = func() error {
			return processGinServerOutput(CLI.Gin.Src + "/apigen.go")
		}
	case "gocli <src> <dst>":
		log.Printf("Gen GO Client")
		src = CLI.Gocli.Src
		processor = func() error {
			return processGoClientOutput(CLI.Gocli.Dst)
		}
	case "pycli <src> <dst>":
		log.Printf("Gen Python Client")
		src = CLI.Pycli.Src
		processor = func() error {
			return processPyClientOutput(CLI.Pycli.Dst)
		}
	case "ts <src> <dst>":
		log.Printf("Gen TS Client")
		src = CLI.Ts.Src
		processor = func() error {
			return processTSClientOutput(CLI.Ts.Dst)
		}
	case "http <src> <dst>":
		log.Printf("Gen Http Client")
		src = CLI.Http.Src
		processor = func() error {
			return processHttpCallOut(CLI.Http.Dst)
		}
	default:
		err = errors.New("unknown option")
	}

	if err != nil {
		panic(err)
	}
	err = load(src)

	if err != nil {
		panic(err)
	}
	err = processor()

	if err != nil {
		panic(err)
	}
}
