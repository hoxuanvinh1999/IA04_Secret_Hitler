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
	c_to_agent := make(map[string]chan voteRequest)
	//Création des agents joueurs
	for _, p := range g.players {
		newChan := make(chan voteRequest)
		//c_to_agent = append(c_to_agent, newChan)
		c_to_agent[p.name] = newChan
		joueur := NewAgentPlayer(p.name, c, newChan, p.role, true, Liberal)
		joueur.Start(g.players)
	}
	g.c_to_agent = c_to_agent
	//fmt.Print(g.c_to_agent)
	time.Sleep(10 * time.Minute)
}

//func (g *game) investigationAvailable() bool {
//	return g.fascistPolicies >= 3
//}
