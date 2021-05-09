package parser

import (
	"fmt"

	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
)

var teamwins int = 0
var enemywins int = 0
var members [5]string

func (p *DemoParser) ProcessResults() {
	fmt.Println("---------------------")
	res := p.isEBCWin()
	var result string = "LOSS"
	if res {
		result = "WIN"
	}
	if teamwins < 16 {
		fmt.Printf("Scoreboard: %d-%d %s\n", enemywins, teamwins, result)
	} else {
		fmt.Printf("Scoreboard: %d-%d %s\n", teamwins, enemywins, result)
	}
	fmt.Print("Team Members:")
	fmt.Println(members)
	fmt.Println("---------------------")
}

// Team 2T 3CT
func (p *DemoParser) isEBCWin() bool {
	var team common.Team = common.Team(1)
	var idx int = 0
	for _, pl := range p.Match.Players {
		if _, found := TEAMMAP[pl.SteamID]; found {
			team = pl.Team.State.Team()
			members[idx] = pl.Name
			idx++
		}
	}

	for _, rnds := range p.Match.Rounds {
		if team == rnds.Winner.State.Team() {
			teamwins++
		} else {
			enemywins++
		}
	}

	if teamwins == 16 {
		return true
	} else if teamwins < 16 {
		return false
	} else {
		return true
	}
}
