package main

import (
	"flag"

	"github.com/zacierka/csgodemoparser/internal/parser"
)

var path string

func init() {
	flag.StringVar(&path, "p", "example/matches/match730_003478395714113896624_1370236049_121.dem", "Path to demo")
}

func main() {
	parser.ParseDemo(path)
}
