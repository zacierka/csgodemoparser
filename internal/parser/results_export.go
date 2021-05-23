package parser

import (
	"fmt"
	"log"

	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var teamwins int = 0
var enemywins int = 0
var members [5]string
var db *gorm.DB = connect()

type Match struct {
	ID          int
	MatchID     string
	Map         string
	ScoreT      int
	ScoreE      int
	Result      uint8 // 0 = loss, 1 = win, 2 = tie
	Playerone   string
	Playertwo   string
	Playerthree string
	Playerfour  string
	Playerfive  string
}

func (p *DemoParser) ProcessResults() {
	setupDatabase(db)
	res := p.isEBCWin()

	match := Match{
		MatchID:     p.Match.ID,
		Map:         p.Match.Map,
		ScoreT:      teamwins,
		ScoreE:      enemywins,
		Result:      res, // loss
		Playerone:   normalizePlayer(p.Match.Players[1].SteamID),
		Playertwo:   normalizePlayer(p.Match.Players[1].SteamID),
		Playerthree: normalizePlayer(p.Match.Players[1].SteamID),
		Playerfour:  normalizePlayer(p.Match.Players[1].SteamID),
		Playerfive:  normalizePlayer(p.Match.Players[1].SteamID),
	}
	// fmt.Println(sort.StringSlice(members[:]))

	id := insertMatch(&match)
	log.Printf("\tresults_export::ProcessResults() inserted match with id %d\n", id)
}

// Team 2T 3CT
func (p *DemoParser) isEBCWin() uint8 {
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
		return 1
	} else if teamwins < 16 {
		return 0
	} else {
		return 2
	}
}

func connect() *gorm.DB {
	dsn := "SQL/PATH/HERE"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func setupDatabase(db *gorm.DB) {
	db.AutoMigrate(&Match{})
}

func insertMatch(match *Match) int {
	db.Create(&match)
	return match.ID
}

func normalizePlayer(steamid uint64) string {
	val := fmt.Sprint(steamid)
	return val
}
