package main

import (
	"time"
)

func main() {

	// go func() {

	// }()

	// go func() {

	// }()

	c := make(chan voteRequest)
	g := newGame([]string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"})
	g.c = c

	g.start()
	c_to_agent := make(map[string]chan voteRequest)
	//Création des agents joueurs
	for _, p := range g.players {
		newChan := make(chan voteRequest)
		c_to_agent[p.name] = newChan
		joueur := NewAgentPlayer(p.name, c, newChan, p.role, true, Liberal, g)
		if joueur.role == Fascist || joueur.role == Hitler {
			for i := 0; i < len(g.players); i++ {
				if g.players[i].role == Fascist || g.players[i].role == Hitler {
					joueur.BothFascists(g.players[i])
				}
			}
		}

		joueur.Start(g.players)
	}
	g.c_to_agent = c_to_agent
	time.Sleep(10 * time.Minute)

	time.Sleep(10 * time.Minute)
}
