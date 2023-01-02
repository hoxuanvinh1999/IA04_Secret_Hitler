package main

import (
	"fmt"
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

func (g *game) voteOnChancellor(president, chancellor player) bool {
	// Vote d'approbation pour le président ("Ja" ou "Nein")
	return true
}

func (g *game) selectPresident() player {
	nextPresident := player{}
	if g.currentPresident == "" {
		nextPresident = g.players[0]
	} else {
		for i, p := range g.players {
			if p.name == g.currentPresident {
				nextPresident = g.players[(i+1)%len(g.players)]
				break
			}
		}
	}
	g.currentPresident = nextPresident.name
	return nextPresident

}

func (g *game) presidentDiscards(president player, cards []string) ([]string, []string) {
	//discarded := make([]string, 0)

	fmt.Printf("%s, choisis une carte à défausser : ", president.name)
	for _, card := range cards {
		fmt.Printf(" %s", card)
	}
	fmt.Println()

	// Scan choix du president
	var choice1 string
	fmt.Scanln(&choice1)

	// On defausse la carte
	g.discard = append(g.discard, choice1)

	// On enleve la carte des cartes à choisir
	cards = remove(cards, choice1)

	return g.discard, cards
}

func (g *game) governmentWasVotedOut() bool {
	// Verifie si plus de la moitié des joueurs a rejeter le gouvernement
	votes := 0

	for _, p := range g.players {
		p.vote = "Rejected"
	}

	for _, p := range g.players {
		if p.vote == "Rejected" {
			votes++
		}
	}
	return votes > len(g.players)/2
}

func (g *game) reshuffle() {
	// Mélange la défausse
	rand.Shuffle(len(g.discard), func(i, j int) {
		g.discard[i], g.discard[j] = g.discard[j], g.discard[i]
	})

	// Ajoute au deck la défausse mélangée
	g.deck = append(g.deck, g.discard...)
	g.discard = make([]string, 0)
}

func (g *game) chancellorEnacts(chancellor player, cards, discarded []string) (string, string) {
	fmt.Printf("%s, choisis une loi à adopter :", chancellor.name)
	for _, card := range cards {
		fmt.Printf(" %s", card)
	}
	fmt.Println()
	var choice string
	var not_choose string
	fmt.Scanln(&choice)
	if choice == cards[0] {
		not_choose = cards[1]
	} else {
		not_choose = cards[0]
	}

	return choice, not_choose
}

func (g *game) printResult() {
	if g.liberalVictory() {
		if !g.hitlerIsAlive() {
			fmt.Println("Hitler est mort, les libéraux ont gagné ! ")
		} else {
			fmt.Println("5 lois libérales ont été votées, les libéraux ont gagné ! ")
		}
	} else if g.fascistVictory() {
		if g.hitlerWasElected() {
			fmt.Printf("%s, qui était Hitler, a été élu. Les fascistes ont gagné !", g.Hitler)
			fmt.Println()
		} else {
			fmt.Println("6 lois fascistes ont été votées, les fascistes ont gagné !")
		}
	} else {
		fmt.Println("Bug/égalité ?")
	}

	fmt.Println("Score final")
	fmt.Println("Lois libérales : ", g.liberalPolicies)
	fmt.Println("Lois fascistes : ", g.fascistPolicies)
}

func (g *game) selectChancellor(president player) player {
	var choice player
	fmt.Printf("%s, choisis un chancelier", president.name)
	for _, p := range g.players {
		if p.name != g.currentChancellor && p.alive && p.name != g.currentPresident {
			fmt.Printf(" %s", p.name)
			choice = p
			break
		}
	}
	fmt.Println()

	return choice
}

// Enleve elements d'un slice
func remove(slice []string, elems ...string) []string {
	for _, elem := range elems {
		for i, e := range slice {
			if e == elem {
				slice = append(slice[:i], slice[i+1:]...)
				break
			}
		}
	}
	return slice
}

func (g *game) start() { //ag *agentMJ
	// Qui est Hitler
	go func() {
		req := <-g.c
		fmt.Printf("agent %q has received %q from %q %q\n", g.ID,
			req.req, req.senderID, req.typerequest)

		for _, p := range g.players {
			if p.role == Hitler {
				g.Hitler = p.name
				break
			}
		}
		//Création de l'agent MJ
		//c := make(chan voteRequest)
		//MJ := NewPongAgent("MJ", c)
		//MJ.Start()

		//Création des agents joueurs
		// for _, p := range g.players {
		// 	joueur := NewAgentPlayer(p.name, c, p.role, true, Liberal)
		// 	joueur.Start()
		// }

		// name := []string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"}
		// role := []string{Liberal, Fascist, Liberal, Hitler, Liberal}

		// for i := 0; i < 5; i++ {
		// 	//id := fmt.Sprintf("pinger n°%d", i)
		// 	pinger := NewAgentPlayer(name[i], c, role[i], true, Liberal)
		// 	pinger.Start()
		// }

		// for _, p := range g.players {
		// 	ag.cout <- Request{"role", "MJ", p.role, ag.cin}
		// }

		// c := make(chan string)
		// req := <- c
		g.currentPresident = "Pierre"
		for !g.isGameOver() {

			// ag.cout <- playerRequestToMj{}

			// ag.cout <- voteRequest{"vote", ag.name, PingString, ag.cin}
			// answer := <-ag.cin
			// fmt.Printf("agent %q has received: %q\n", ag.name, answer)

			//req := <-ag.c
			//answer := <-ag.cin
			//fmt.Printf("le test : agent %q has received %q from %q %q\n", ag.ID,
			//	req.req, req.senderID, req.typerequest)
			//roles := make([]string, numPlayers)
			//newreq := make(chan, voteRequest)

			newreq := <-g.c
			fmt.Printf("agent %q has received %q from %q %q\n", g.ID,
				req.req, req.senderID, req.typerequest)

			for _, p := range g.players {
				fmt.Print("p.name :", p.name)
				newreq = <-g.c
				fmt.Println("Curreent pres : ", g.currentPresident)
				fmt.Println("Sender ID :", newreq.senderID)
				if g.currentPresident == newreq.senderID {
					break
				}
				fmt.Printf("agent %q has received %q from %q %q\n", g.ID,
					newreq.req, newreq.senderID, newreq.typerequest)
				time.Sleep(1 * time.Second)
				fmt.Println("newreq sendeid", newreq.senderID)
			}

			time.Sleep(1 * time.Second)
			fmt.Print("newreq ?", newreq.senderID)
			//lui demander de choisir un president
			go g.handlePing(newreq)

			time.Sleep(5 * time.Second)
			// Choix du president pour ce tour
			president := g.selectPresident()
			g.currentPresident = president.name
			// Le president propose un chancelier
			chancellor := g.selectChancellor(president)

			// Si le vote passe, on fait le tour, sinon tour suivant
			if g.voteOnChancellor(president, chancellor) {
				g.currentChancellor = chancellor.name
				// Le president pioche 3 cartes et en defausse une
				if g.hitlerIsAlive() && g.hitlerWasElected() && g.fascistPolicies >= 3 {
					break
				}
				cards := g.drawCards(3)
				discarded, cards := g.presidentDiscards(president, cards)
				// Le chancelier choisit une des deux cartes et defausse l'autre
				enacted, not_enacted := g.chancellorEnacts(chancellor, cards, discarded)
				// On defausse la carte non choisit
				g.discard = append(g.discard, not_enacted)
				// On ajoute la carte choisit
				g.enactPolicy(enacted)
			}

			if g.isGameOver() {
				break
			}
			// Si le gouvernement est refuse, on choisit un nouveau president
			if g.governmentWasVotedOut() {
				fmt.Println("gov voted out")
				president = g.selectPresident()
			}

		}

		g.printResult()
	}()
}
