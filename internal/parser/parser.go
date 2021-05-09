package parser

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	demoinfocs "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

type DemoParser struct {
	parser        demoinfocs.Parser
	Match         *MatchData
	CurrentRound  byte
	RoundStart    time.Duration
	RoundOngoing  bool
	SidesSwitched bool
}

// MatchData holds information about the match itself.
type MatchData struct {
	ID       string
	Map      string
	Header   *common.DemoHeader
	Players  []*Player
	Teams    [2]*Team
	Duration time.Duration
	Time     time.Time
	Rounds   []*Round
}

// Team represents a team and links to it's players.
type Team struct {
	StartedAs common.Team
	State     *common.TeamState
	Players   []*Player
}

// Player represents one player either as T or CT.
type Player struct {
	SteamID uint64
	Name    string
	Team    *Team
}

// Round contains information about one round.
type Round struct {
	Duration  time.Duration
	Kills     []*Kill
	Winner    *Team
	WinReason events.RoundEndReason
	MVP       *Player
}

// Kill holds information about a kill that happenend during the match.
type Kill struct {
	Time          time.Duration
	Victim        *Player
	Killer        *Player
	Assister      *Player
	Weapon        common.EquipmentType
	IsDuringRound bool
	IsHeadshot    bool
	AssistedFlash bool
	AttackerBlind bool
	NoScope       bool
	ThroughSmoke  bool
	ThroughWall   bool
}

func (p *DemoParser) ParseDemo(path string) error {
	match_ID := strings.Split(path, "_")[1]
	p.Match = &MatchData{ID: match_ID}
	f, err := os.Open(path)

	if err != nil {
		return err
	}
	log.Println("\tParsing:", path)

	p.parser = demoinfocs.NewParser(f)

	header, _ := p.parser.ParseHeader()
	p.Match.Header = &header

	// Register all handler
	p.parser.RegisterEventHandler(p.handleMatchStart)
	p.parser.RegisterEventHandler(p.handleGamePhaseChanged)
	p.parser.RegisterEventHandler(p.handleKill)
	p.parser.RegisterEventHandler(p.handleRoundStart)
	p.parser.RegisterEventHandler(p.handleRoundEnd)

	defer p.parser.Close()
	defer f.Close()

	//var res bool = p.isEBCWin()
	//fmt.Println(res)

	return p.parser.ParseToEnd()
}

func (p *DemoParser) getPlayer(player *common.Player) (*Player, error) {
	if player.IsBot {
		return nil, errors.New("Player is a bot")
	}

	for _, localPlayer := range p.Match.Players {
		if player.SteamID64 == localPlayer.SteamID {
			return localPlayer, nil
		}
	}

	for _, gamePlayer := range p.parser.GameState().Participants().Playing() {
		if player.SteamID64 == gamePlayer.SteamID64 {
			return p.AddPlayer(player), nil
		}
	}

	return nil, errors.New("Player not found in local match struct " + strconv.FormatUint(player.SteamID64, 10))
}

// GetTeamIndex returns 0 for T, 1 for CT and 2 for everything else.
func GetTeamIndex(team common.Team, sidesSwitched bool) byte {
	if team == common.TeamTerrorists {
		if !sidesSwitched {
			return 0
		}
		return 1
	} else if team == common.TeamCounterTerrorists {
		if !sidesSwitched {
			return 1
		}
		return 0
	}

	// Could also return an error here but we do not expect this to happen.
	return 2
}

// AddPlayer adds a player to the game and returns the pointer.
func (p *DemoParser) AddPlayer(player *common.Player) *Player {
	teamID := GetTeamIndex(player.Team, p.SidesSwitched)
	teams := p.Match.Teams
	teamPlayers := teams[teamID].Players

	customPlayer := &Player{SteamID: player.SteamID64, Name: player.Name, Team: teams[teamID]}

	teams[teamID].Players = append(teamPlayers, customPlayer)
	p.Match.Players = append(p.Match.Players, customPlayer)
	return customPlayer
}
