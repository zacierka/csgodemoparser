package parser

import "fmt"

func (m *MatchData) ProcessResults() {
	fmt.Printf("Total Players: %d\n", len(m.Players))
}
