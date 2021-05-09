package main

import (
	"flag"
	"log"

	"github.com/zacierka/csgodemoparser/internal/parser"
)

var path string

func init() {
	flag.StringVar(&path, "p", "example/matches/match730_003478395714113896624_1370236049_121.dem", "Path to demo")
}

func main() {
	parser := &parser.DemoParser{}
	err := parser.ParseDemo(path)

	if err != nil {
		panic("Error Parsing File")
	}
	log.Println("\tFinished Parsing:", path)

	log.Printf("\tProcessing Results for game ID: %s", parser.Match.ID)
	parser.Match.ProcessResults()
}
