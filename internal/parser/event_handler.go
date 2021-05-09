package parser

import (
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

// Inits the players and teams.
func (p *DemoParser) handleMatchStart(e events.MatchStart) {
	p.Match.Map = p.Match.Header.MapName
	p.SidesSwitched = false

	gameState := p.parser.GameState()

	// Create teams.
	ct := gameState.TeamCounterTerrorists()
	t := gameState.TeamTerrorists()

	p.Match.Teams[GetTeamIndex(t.Team(), p.SidesSwitched)] = &Team{State: t, StartedAs: common.TeamTerrorists}
	p.Match.Teams[GetTeamIndex(ct.Team(), p.SidesSwitched)] = &Team{State: ct, StartedAs: common.TeamCounterTerrorists}

	// Create players and map them to the teams.
	for _, player := range gameState.Participants().Playing() {
		if player.IsBot {
			continue
		}

		p.AddPlayer(player)
	}
}

func (p *DemoParser) handleGamePhaseChanged(e events.GamePhaseChanged) {
	switch e.NewGamePhase {
	case common.GamePhaseInit:
		p.SidesSwitched = false
	case common.GamePhaseTeamSideSwitch:
		p.SidesSwitched = !p.SidesSwitched
	case common.GamePhaseGameEnded:
		p.Match.Duration = p.parser.CurrentTime()
	}
}

func (p *DemoParser) handleRoundStart(e events.RoundStart) {
	if p.RoundOngoing {
		return
	}
	p.CurrentRound++
	p.RoundOngoing = true
	p.RoundStart = p.parser.CurrentTime()
	p.Match.Rounds = append(p.Match.Rounds, &Round{})
}

func (p *DemoParser) handleRoundEnd(e events.RoundEnd) {
	if !p.RoundOngoing {
		return
	}

	p.RoundOngoing = false
	round := p.Match.Rounds[p.CurrentRound-1]

	round.Winner = p.Match.Teams[GetTeamIndex(e.Winner, p.SidesSwitched)]
	round.WinReason = e.Reason
	round.Duration = p.parser.CurrentTime() - p.RoundStart
}

func (p *DemoParser) handleKill(e events.Kill) {
	if p.parser.GameState().IsWarmupPeriod() || p.CurrentRound == 0 {
		return
	}

	round := p.Match.Rounds[p.CurrentRound-1]
	kill := &Kill{Time: p.parser.CurrentTime(), Weapon: e.Weapon.Type, IsHeadshot: e.IsHeadshot,
		AssistedFlash: e.AssistedFlash, AttackerBlind: e.AttackerBlind, NoScope: e.NoScope,
		ThroughSmoke: e.ThroughSmoke, ThroughWall: e.IsWallBang(), IsDuringRound: p.RoundOngoing}

	victim, err := p.getPlayer(e.Victim)
	if err == nil {
		kill.Victim = victim
	}

	// Add optional killer if player died e.g. through fall damage.
	if e.Killer != nil {
		killer, err := p.getPlayer(e.Killer)
		if err == nil {
			kill.Killer = killer
		}
	}

	// Add optional assister.
	if e.Assister != nil {
		assister, err := p.getPlayer(e.Assister)
		if err == nil {
			kill.Assister = assister
		}
	}

	round.Kills = append(round.Kills, kill)
}
