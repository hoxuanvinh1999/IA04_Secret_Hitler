package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const (
	Liberal = "Liberal"
	Fascist = "Fascist"
	Hitler  = "Hitler"
)

const (
	numberOfPlayers  = 10
	numberOfFascists = 3
	numberOfLiberals = 7
	numberOfCards    = 17
)
const (
	President  = "President"
	Chancellor = "Chancellor"
)

type player struct {
	name  string
	role  string
	alive bool
	vote  string
}

type game struct {
	players                  []player
	deck                     []string
	discard                  []string
	mu                       sync.Mutex
	logs                     []string
	liberalPolicies          int
	fascistPolicies          int
	investigationAvailable   bool
	specialElectionAvailable bool
	executionAvailable       bool
	currentPresident         string
	currentChancellor        string
	Hitler                   string
}

func (g *game) log(format string, a ...interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.logs = append(g.logs, fmt.Sprintf(format, a...))
}

func (g *game) start() {
	// Qui est Hitler
	for _, p := range g.players {
		if p.role == Hitler {
			g.Hitler = p.name
			break
		}
	}

	for !g.isGameOver() {
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
}

func main() {
	g := newGame([]string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"})
	g.start()
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

//func (g *game) investigationAvailable() bool {
//	return g.fascistPolicies >= 3
//}
