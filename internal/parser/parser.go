package parser

import (
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
)

func ParseDemo(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	// Register Events
	p.RegisterEventHandler(HandlePlayerConnect)
	// p.RegisterEventHandler(HandleKill)

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}
}
