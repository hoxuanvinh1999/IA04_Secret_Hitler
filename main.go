package main

import (
	"time"
)

func main() {
	c := make(chan voteRequest)
	g := newGame([]string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"})
	g.c = c
	// c := make(chan voteRequest)
	// ponger := NewPongAgent("ponger", c)

	// ponger.Start()
	// name := []string{"Vinh", "Wassim", "Pierre", "Sylvain", "Nathan"}
	// role := []string{Liberal, Fascist, Liberal, Hitler, Liberal}
	// for i := 0; i < 5; i++ {
	// 	//id := fmt.Sprintf("pinger n°%d", i)
	// 	pinger := NewAgentPlayer(name[i], c, role[i], true, Liberal)
	// 	pinger.Start()
	// }
	// time.Sleep(time.Second)

	//g.start(ponger)
	g.start()

	//Création des agents joueurs
	for _, p := range g.players {
		joueur := NewAgentPlayer(p.name, c, p.role, true, Liberal)
		joueur.Start()
	}

	time.Sleep(10 * time.Minute)
}

//func (g *game) investigationAvailable() bool {
//	return g.fascistPolicies >= 3
//}
