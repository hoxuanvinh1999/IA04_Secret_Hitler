package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func newGame(names []string) *game {

	rand.Seed(time.Now().UnixNano())
	numPlayers := len(names)
	roles := make([]string, numPlayers)
	for i := 0; i < numPlayers/2+1; i++ {
		roles[i] = Liberal
	}
	roles[numPlayers/2+1] = Hitler
	for i := numPlayers/2 + 2; i < numPlayers; i++ {
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

	result := &game{
		players: players,
		deck:    deck,
		discard: make([]string, 0),
		logs:    make([]string, 0),
	}
	return result
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
		if p.name == g.currentChancellor.name && p.role == Hitler {
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
		fmt.Print("Une loi libérale a été votée.\n")
	} else if policy == Fascist {
		g.fascistPolicies++
		fmt.Print("Une loi fasciste a été votée.\n")
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
	nb_Ja := 0
	nb_Nein := 0
	for _, p := range g.players {
		fmt.Printf("%s, vote Ja ou Nein pour élire : %s \n", p.name, chancellor.name)
		g.c_to_agent[p.name] <- voteRequest{"vote", "MJ", "MJ", PingString, g.c, chancellor, []string{"Liberal"}, true, game_vide}
		answer := <-g.c
		if answer.Ja {
			nb_Ja += 1
		} else {
			nb_Nein += 1
		}

		time.Sleep(100 * time.Millisecond)
	}

	if nb_Ja > nb_Nein {
		fmt.Println("Le résultat de l'élection est JA !")
		g.result_vote = true
		g.nombre_echec = 0
	} else {
		fmt.Println("Le résultat de l'élection est NEIN !")
		g.result_vote = false
		g.nombre_echec++
		if g.nombre_echec >= 3 {
			g.chaos = true
			g.nombre_echec = 0
		}

	}

	return nb_Ja > nb_Nein
}

func (g *game) selectPresident() player {
	nextPresident := player{}
	if g.currentPresident.name == "" {
		nextPresident = g.players[0]
	} else {
		for i, p := range g.players {
			if p.name == g.currentPresident.name {
				nextPresident = g.players[(i+1)%len(g.players)]
				break
			}
		}
	}

	// go func() {
	// 	website := Website{"Secret Hitler", time.Now().Format(time.Stamp),
	// 		g.website.Players_name,
	// 		g.website.Players_side,
	// 		g.website.Players_alive,
	// 		[]string{"_", "_", "_", "_", "_"},
	// 		[]string{"_", "_", "_", "_", "_", "_"},
	// 		g.website.Current_President,
	// 		g.website.Current_President,
	// 	}
	// 	template := template.Must(template.ParseFiles("web/main.html"))
	// 	go func() {
	// 		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// 		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 			if game_title := r.FormValue("white_background"); game_title != "" {
	// 				website.Game_title = game_title
	// 			}
	// 			if err := template.ExecuteTemplate(w, "main.html", website); err != nil {
	// 				http.Error(w, err.Error(), http.StatusInternalServerError)
	// 			}
	// 		})
	// 		fmt.Println(http.ListenAndServe(":8001", nil))
	// 	}()

	// }()

	g.currentPresident = nextPresident
	return nextPresident

}

// func (g *game) selectChancellor(president player) player {
// 	var choice player
// 	fmt.Printf("%s, choisis un chancelier", president.name)
// 	for _, p := range g.players {
// 		if p.name != g.currentChancellor && p.alive && p.name != g.currentPresident {
// 			fmt.Printf(" %s", p.name)
// 			choice = p
// 			break
// 		}
// 	}
// 	fmt.Printf("%s, propose pour chancelier", choice.name)
// 	fmt.Println()

// 	return choice
// }

func (g *game) selectChancellor(president player, chancelier player) player {
	var choice player
	fmt.Printf("%s, choisis un chancelier\n", president.name)
	choice = chancelier
	fmt.Printf("%s, propose pour chancelier %s\n", president.name, choice.name)
	g.propChancellor = choice
	return choice
}

// func (g *game) selectChancellor(presi player, p player) player {
// 	nextPresident := player{}
// 	nextPresident = p
// 	g.currentPresident = nextPresident.name
// 	fmt.Println("Le chancelier %q choisit pour président %q.\n", ,nextPresident)
// 	return nextPresident

// }

// func (g *game) presidentDiscards(president player, cards []string) ([]string, []string) {
// 	//discarded := make([]string, 0)
// 	fmt.Printf("%s, choisis une carte à défausser : ", president.name)
// 	for _, card := range cards {
// 		fmt.Printf(" %s", card)
// 	}
// 	fmt.Println()
// 	// Scan choix du president
// 	var choice1 string
// 	fmt.Scanln(&choice1)
// 	// On defausse la carte
// 	g.discard = append(g.discard, choice1)
// 	// On enleve la carte des cartes à choisir
// 	cards = remove(cards, choice1)

// 	return g.discard, cards
// }

func (g *game) presidentDiscards(president player, cards []string) ([]string, []string) {
	var choice string
	for _, p := range g.players {
		if g.currentPresident.name == p.name {
			fmt.Printf("%s, choisis une carte à défausser : ", president.name)
			for _, card := range cards {
				fmt.Printf(" %s", card)
			}
			fmt.Println()

			g.c_to_agent[p.name] <- voteRequest{"choisisdiscards", "MJ", "MJ", PingString, g.c, p, cards, true, game_vide}
			answer := <-g.c

			choice = answer.cards[2]
			cards = answer.cards[0:2]
		}
		time.Sleep(200 * time.Millisecond)
	}

	// On defausse la carte
	g.discard = append(g.discard, choice)

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

// func (g *game) chancellorEnacts(chancellor player, cards, discarded []string) (string, string) {
// 	fmt.Printf("%s, choisis une loi à adopter :", chancellor.name)
// 	for _, card := range cards {
// 		fmt.Printf(" %s", card)
// 	}
// 	fmt.Println()
// 	var choice string
// 	var not_choose string
// 	fmt.Scanln(&choice)
// 	if choice == cards[0] {
// 		not_choose = cards[1]
// 	} else {
// 		not_choose = cards[0]
// 	}

// 	return choice, not_choose
// }

func (g *game) chancellorEnacts(chancellor player, cards, discarded []string) (string, string) {
	var choice string
	var not_choose string
	for _, p := range g.players {
		if chancellor.name == p.name {
			fmt.Printf("%s, choisis une loi à adopter :", chancellor.name)
			for _, card := range cards {
				fmt.Printf(" %s", card)
			}
			fmt.Println()

			g.c_to_agent[p.name] <- voteRequest{"enact", "MJ", "MJ", PingString, g.c, p, cards, true, game_vide}
			answer := <-g.c

			choice = answer.cards[1]
			not_choose = answer.cards[0]
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Carte posée :", choice)
	g.prevPresident = g.currentPresident
	g.prevChancellor = g.currentChancellor

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

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	go func() {
		http.HandleFunc("/websoc", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Print("upgrade:", err)
			}
			defer conn.Close()
			for i := 0; i < len(g.players); i++ {
				time.Sleep(200 * time.Millisecond)
				message := "0" + strconv.Itoa(i) + g.players[i].name
				conn.WriteMessage(websocket.TextMessage, []byte(message))
			}
			for i := 0; i < len(g.players); i++ {
				time.Sleep(200 * time.Millisecond)
				message := "1" + strconv.Itoa(i) + g.players[i].role
				conn.WriteMessage(websocket.TextMessage, []byte(message))
			}
			for i := 0; i < len(g.players); i++ {
				time.Sleep(200 * time.Millisecond)
				message := "2" + strconv.Itoa(i) + "alive"
				conn.WriteMessage(websocket.TextMessage, []byte(message))
			}

			time.Sleep(200 * time.Millisecond)
			message_candidat := "7" + g.propChancellor.name
			conn.WriteMessage(websocket.TextMessage, []byte(message_candidat))
			for i := 0; i < 1000; i++ {
				time.Sleep(200 * time.Millisecond)
				message_candidat := "7" + g.propChancellor.name
				conn.WriteMessage(websocket.TextMessage, []byte(message_candidat))
				time.Sleep(200 * time.Millisecond)
				if g.result_vote {
					message_result := "8" + "Ja"
					conn.WriteMessage(websocket.TextMessage, []byte(message_result))
				} else {
					message_result := "8" + "Nein"
					conn.WriteMessage(websocket.TextMessage, []byte(message_result))
				}
				time.Sleep(200 * time.Millisecond)
				message_president := "3" + g.currentPresident.name
				conn.WriteMessage(websocket.TextMessage, []byte(message_president))
				message_chancellor := "4" + g.currentChancellor.name
				conn.WriteMessage(websocket.TextMessage, []byte(message_chancellor))
				time.Sleep(200 * time.Millisecond)
				if g.liberalPolicies > 0 {
					message_board_liberal := "5" + strconv.Itoa(g.liberalPolicies)
					conn.WriteMessage(websocket.TextMessage, []byte(message_board_liberal))
				}
				time.Sleep(200 * time.Millisecond)
				if g.fascistPolicies > 0 {
					message_board_fascist := "6" + strconv.Itoa(g.fascistPolicies)
					conn.WriteMessage(websocket.TextMessage, []byte(message_board_fascist))
				}
				// for i := 0; i < len(g.players); i++ {
				// 	time.Sleep(200 * time.Millisecond)
				// 	if !g.players[i].alive {
				// 		message := "2" + strconv.Itoa(i) + "dead"
				// 		conn.WriteMessage(websocket.TextMessage, []byte(message))
				// 	}
				// }
			}
		})
		fmt.Println("server running on port 8000")
		log.Fatal(http.ListenAndServe(":8000", nil))
	}()

	// Qui est Hitler
	go func() {
		for _, p := range g.players {
			//give time for the web to load
			time.Sleep(1 * time.Second)
			fmt.Printf("%q est %q \n", p.name, p.role)
		}

		//req := <-g.c
		//fmt.Printf("Le mj a recu : je suis %q, mon role est %q j'envoie %q, et j'envoie une requête de type %q \n", req.senderID, req.role, req.req, req.typerequest)
		for _, p := range g.players {
			if p.role == Hitler {
				g.Hitler = p.name
				break
			}
		}

		for !g.isGameOver() {

			// Choix du president pour ce tour
			president := g.selectPresident()

			for _, p := range g.players {
				if g.currentPresident.name == p.name {
					g.c_to_agent[p.name] <- voteRequest{"choisischancelier", "MJ", "MJ", PingString, g.c, p, []string{}, true, game_vide}

					//fmt.Printf("%q propose pour chancelier %q\n", newreq.senderID, newreq.playerpres.name)
				}
				time.Sleep(200 * time.Millisecond)
			}

			newreq := <-g.c

			time.Sleep(200 * time.Millisecond)

			g.currentPresident.name = president.name
			// Le president propose un chancelier
			chancellor := g.selectChancellor(president, newreq.playerpres)

			// Si le vote passe, on fait le tour, sinon tour suivant
			if g.voteOnChancellor(president, chancellor) {
				g.currentChancellor = chancellor
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
			if g.chaos {
				fmt.Print("C'est le chaos !!")
				card := g.drawCards(1)
				g.enactPolicy(card[0])
				g.chaos = false

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
		g.end = true
		time.Sleep(30 * time.Second)
	}()
}
