package main

import (
	"fmt"
)

func (g *game) log(format string, a ...interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.logs = append(g.logs, fmt.Sprintf(format, a...))
}

/*
type agentPlayer struct {
	name  string
	role  string
	alive bool
	vote  string
	cin   chan string
	cout  chan Request
}

func NewAgentPlayer(name string, cout chan Request, role string, alive bool, vote string) *agentPlayer {
	cin := make(chan string)
	return &agentPlayer{name, role, alive, vote, cin, cout}
}

func (ag *agentPlayer) Start() {
	go func() {
		for {
			ag.cout <- Request{ag.name, PingString, ag.cin}
			answer := <-ag.cin
			fmt.Printf("agent %q has received: %q\n", ag.name, answer)
			time.Sleep(time.Second)
		}
	}()
}
*/

func (ag1 *agentPlayer) BeliefUp(ag2 player) {
	for i := 0; i < len(ag1.currentGame.players); i++ {
		if ag1.currentGame.players[i] == ag2 {
			ag1.beliefs[ag1.currentGame.players[i]] += 1
		}
	}
}

func (ag1 *agentPlayer) BeliefDown(ag2 player) {
	for i := 0; i < len(ag1.currentGame.players); i++ {
		if ag1.currentGame.players[i] == ag2 {
			ag1.beliefs[ag1.currentGame.players[i]] -= 1
		}
	}
}

func (ag1 *agentPlayer) BothFascists(ag2 player) {
	for i := 0; i < len(ag1.currentGame.players); i++ {
		if ag1.currentGame.players[i] == ag2 {
			ag1.beliefs[ag1.currentGame.players[i]] = 1
		}
	}
}
