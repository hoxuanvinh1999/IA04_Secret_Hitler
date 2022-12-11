package main

import (
	"math/rand"
	"time"
)

func newGame(names []string) *game {
	rand.Seed(time.Now().UnixNano())
	numPlayers := len(names)
	roles := make([]string, numPlayers)
	for i := 0; i < numPlayers/2; i++ {
		roles[i] = Liberal
	}
	roles[numPlayers/2] = Hitler
	for i := numPlayers/2 + 1; i < numPlayers; i++ {
		roles[i] = Fascist
	}

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	players := make([]player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = player{
			name:  names[i],
			role:  roles[i],
			alive: true,
		}
	}

	deck := make([]string, 17)
	for i := 0; i < 6; i++ {
		deck[i] = Liberal
	}
	for i := 6; i < 17; i++ {
		deck[i] = Fascist
	}
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return &game{
		players: players,
		deck:    deck,
		discard: make([]string, 0),
		logs:    make([]string, 0),
	}
}

func (g *game) isGameOver() bool {
	return g.liberalVictory() || g.fascistVictory()
}

func (g *game) liberalVictory() bool {
	return (!g.hitlerIsAlive()) || (g.liberalPolicies == 5)
}

func (g *game) hitlerIsAlive() bool {
	hitlerAlive := false
	for _, p := range g.players {
		if p.role == Hitler && p.alive {
			hitlerAlive = true
			break
		}
	}
	return hitlerAlive
}

func (g *game) hitlerWasElected() bool {
	// Vérifie si le chancelier actuel est Hitler
	check := false
	for _, p := range g.players {
		if p.name == g.currentChancellor && p.role == Hitler {
			check = true
		}
	}
	return check
}

func (g *game) fascistVictory() bool {
	return (g.hitlerIsAlive() && g.hitlerWasElected() && g.fascistPolicies >= 3) || (g.fascistPolicies == 6)
}

func (g *game) drawCards(num int) []string {
	cards := make([]string, num)
	for i := 0; i < num; i++ {
		// Remélange le deck si nécessaire
		if len(g.deck) < 3 {
			g.reshuffle()
		}

		// Pioche une carte
		cards[i] = g.deck[0]
		g.deck = g.deck[1:]
	}
	return cards
}

func (g *game) enactPolicy(policy string) {

	if policy == Liberal {
		g.liberalPolicies++
	} else if policy == Fascist {
		g.fascistPolicies++

		// En fonction du nombre de lois fascistes, des actions sont possibles
		if g.fascistPolicies == 3 {
			// Le président peut regarder un rôle
			g.investigationAvailable = true
		} else if g.fascistPolicies == 4 {
			// Président peut déclencher une élection spéciale, et donc choisir le futur candidat à la présidence
			g.specialElectionAvailable = true
		} else if g.fascistPolicies == 5 {
			// Le président peut exécuter un joueur
			g.executionAvailable = true
		}
	}
}
