package parser

import (
	"fmt"

	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

var pcount int

func Init() {
	pcount = 0
}

func HandlePlayerConnect(e events.PlayerConnect) {
	pcount++
	fmt.Printf("%d) %s\n", pcount, e.Player.Name)
}

func HandleKill(e events.Kill) {
	var hs string
	if e.IsHeadshot {
		hs = " (HS)"
	}
	var wallBang string
	if e.PenetratedObjects > 0 {
		wallBang = " (WB)"
	}
	fmt.Printf("%s <%v%s%s> %s\n", e.Killer, e.Weapon, hs, wallBang, e.Victim)
}
