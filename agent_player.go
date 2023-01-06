package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const PingString = "ping"
const PongString = "pong"

var c chan voteRequest
var player_vide = player{"rien", "rien", false, "rien", c}
var choice = ""

var game_vide = newGame([]string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"})

type voteRequest struct {
	typerequest string
	senderID    string
	role        string
	req         string
	c           chan voteRequest
	playerpres  player
	cards       []string
	Ja          bool
	game        *game
	nombre      float64
}

type Agent interface {
	Start()
}

type agentPlayer struct {
	name         string
	role         string
	alive        bool
	vote         string
	cin          chan voteRequest
	cout         chan voteRequest
	beliefs      map[player]int
	currentGame  *game
	menteur      float64
	perspicacite float64
}

//Crée un agentPlayer
func NewAgentPlayer(name string, cout chan voteRequest, cin chan voteRequest, role string, alive bool, vote string, currentGame *game, menteur float64, perspicacite float64) *agentPlayer {
	//cin := make(chan voteRequest)
	beliefs := make(map[player]int, len(currentGame.players))
	for i := 0; i < len(currentGame.players); i++ {
		beliefs[currentGame.players[i]] = 3
	}

	return &agentPlayer{name, role, alive, vote, cin, cout, beliefs, currentGame, menteur, perspicacite}
}

//Augmente la suspicion d'un agent pour un autre
func BeliefUp(ag1 *agentPlayer, ag2 player) {
	ag1.beliefs[ag2] += 1
}

//Diminue la suspicion d'un agent pour un autre
func BeliefDown(ag1 *agentPlayer, ag2 player) {
	ag1.beliefs[ag2] -= 1
}

//Vérifie que les deux sont fascistes
func (ag1 *agentPlayer) BothFascists(ag2 player) {
	for i := 0; i < len(ag1.currentGame.players); i++ {
		if ag1.currentGame.players[i] == ag2 {
			ag1.beliefs[ag1.currentGame.players[i]] = 1
		}
	}
}
//Loi normale
func RandomNormal(mean, stdDev float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return z0*stdDev + mean
}

func (ag *agentPlayer) Start(list_player []player) {
	go func() {
		for {
			//ag.cout <- voteRequest{"vote", ag.name, ag.role, PingString, ag.cin, player_vide}
			answer := <-ag.cin
			//Si l'agent doit choisir un chancelier
			if answer.typerequest == "choisischancelier" {
				fmt.Printf("agent %q has received: %q\n", ag.name, answer.typerequest)
				fmt.Printf("Je vais choisir un chancelier.\n")
				for _, p := range list_player {
					//En fonction de son rôle et de ses croyances il choisit
					if ag.role == Fascist || ag.role == Hitler {
						if (p.role == Fascist || p.role == Hitler) && (ag.name != p.name) && (p.name != answer.game.currentChancellor.name) && (p.alive) {
							player_vide = p
							break
						} else {
							for _, q := range list_player {
								if (ag.name != q.name) && (q.name != answer.game.currentChancellor.name) && (q.alive) {
									player_vide = q
									break
								}
							}

						}
					} else {
						if (ag.beliefs[p] > 2) && (ag.name != p.name) && (p.alive) {
							player_vide = p
						} else {
							for _, q := range list_player {
								if (ag.name != q.name) && (q.name != answer.game.currentChancellor.name) && (q.alive) {
									player_vide = q
									break
								}
							}

						}
					}
				}
				if player_vide == answer.game.currentChancellor {
					for _, p := range list_player {
						if (player_vide != p) && (p.alive) {
							player_vide = p
						}
					}
				}
				//Répond le président qu'il porte pour candidat
				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, []string{}, true, game_vide, 0}

				//answer := <-ag.cin
			//Si l'agent doit choisir une carte à défausser
			} else if answer.typerequest == "choisisdiscards" {
				fmt.Printf("Agent %q a reçu une demande de %q avec pour cartes : %q\n", ag.name, answer.typerequest, answer.cards)
				if ag.role == Fascist || ag.role == Hitler {
					for i := range answer.cards {
						if answer.cards[i] == "Liberal" {
							choice = answer.cards[i]
							break
						} else {
							choice = answer.cards[0]
						}
					}
				} else {
					for i := range answer.cards {
						if answer.cards[i] == "Fascist" {
							choice = answer.cards[i]
							break
						} else {
							choice = answer.cards[0]
						}
					}
				}
				answer.cards = remove(answer.cards, choice)
				answer.cards = append(answer.cards, choice)

				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, answer.cards, true, game_vide, 0}

			//Si l'agent doit choisir une carte à promulguer
			} else if answer.typerequest == "enact" {
				fmt.Printf("Agent %q a reçu une demande de %q avec pour cartes : %q\n", ag.name, answer.typerequest, answer.cards)

				if ag.role == Fascist || ag.role == Hitler {
					for i := range answer.cards {
						if answer.cards[i] == "Liberal" {
							choice = answer.cards[0]
						} else {
							choice = answer.cards[i]
							break
						}
					}
				} else {
					for i := range answer.cards {
						if answer.cards[i] == "Fascist" {
							choice = answer.cards[0]
						} else {
							choice = answer.cards[i]
							break
						}
					}
				}

				fmt.Println("choice", choice)

				answer.cards = remove(answer.cards, choice)
				answer.cards = append(answer.cards, choice)

				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, answer.cards, true, game_vide, 0}

			//Si l'agent reçoit une question
			} else if answer.typerequest == "question" {
				fmt.Print(ag.name, " es-tu fasciste ? \n")

				mean := ag.menteur
				stdDev := 1.0
				rand_float := RandomNormal(mean, stdDev)

				//rand_int := rand.Intn(5) + 1
				reponse := ""
				if rand_float < 1 {
					reponse = "Non, bien sûr, je ne suis pas fasciste. J'agis en tant que libéral depuis le début, et je le suis."
				} else if 1 <= rand_float && rand_float < 2 {
					reponse = "Je ne suis pas fasciste. J'agis en tant que libéral depuis le début."
				} else if 2 <= rand_float && rand_float < 3 {
					reponse = "Je ne suis pas fasciste."
				} else if 3 <= rand_float && rand_float < 4 {
					reponse = "Euh... Non, non, je suis bien libéral"
				} else if 4 <= rand_float && rand_float < 5 {
					reponse = "Euh... Mais quoi... Pourquoi vous me soupçounnez toujours, ce n'est pas juste, j'en ai marre de ce jeu !!"
				}
				ag.cout <- voteRequest{"reponse", ag.name, ag.role, reponse, ag.cin, player_vide, answer.cards[0:1], answer.Ja, game_vide, 0}

			} else if answer.typerequest == "vote" {

				//fmt.Print(ag)
				if (ag.role == Liberal) && (ag.currentGame.currentChancellor.name != ag.name) {
					if choice == Fascist {
						BeliefDown(ag, ag.currentGame.prevPresident)
						BeliefDown(ag, ag.currentGame.prevChancellor)
						//fmt.Println(ag.beliefs)
					} else if choice == Liberal {
						BeliefUp(ag, ag.currentGame.prevPresident)
						BeliefUp(ag, ag.currentGame.prevChancellor)
						//fmt.Println(ag.beliefs)
					}

				}

				fmt.Printf("Agent %q a reçu une demande de %q pour élire %q chancelier. \n", ag.name, answer.typerequest, answer.playerpres.name)

				if ag.role == Fascist || ag.role == Hitler {
					if answer.playerpres.role == Fascist || answer.playerpres.role == Hitler {
						answer.Ja = true
					} else {
						answer.Ja = false
					}
				} else { // le joueur est libéral
					if ag.beliefs[answer.playerpres] < 3 {
						answer.Ja = false
					} else {
						answer.Ja = true
					}
				}

				if answer.Ja {
					fmt.Println("Il vote Ja !")
				} else {
					fmt.Println("Il vote Nein !")
				}

				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, answer.cards[0:1], answer.Ja, game_vide, 0}

			//Si le joueur reçoit une réponse à la question
			} else if answer.typerequest == "reponse" {

				mean := ag.perspicacite
				stdDev := 0.5
				rand_float := RandomNormal(mean, stdDev)
				if rand_float > 0.5 {
					if answer.nombre == 1 {
						BeliefUp(ag, ag.currentGame.propChancellor)
					} else if answer.nombre == 0.5 {
						rand_int2 := rand.Intn(2) + 1
						if rand_int2 == 1 {
							BeliefUp(ag, ag.currentGame.propChancellor)
						}
					} else if answer.nombre == -0.5 {
						rand_int2 := rand.Intn(2) + 1
						if rand_int2 == 1 {
							BeliefDown(ag, ag.currentGame.propChancellor)
						}
					} else if answer.nombre == -1 {
						BeliefDown(ag, ag.currentGame.propChancellor)
					}
				}

			//Si le joueur doit exécuter un autre
			} else if answer.typerequest == "execute" {
				if ag.role == Fascist || ag.role == Hitler {
					minValue := math.MaxInt32
					var agents_execute player
					for agentss, value := range ag.beliefs {
						if (value < minValue) && (agentss.name != ag.name) {
							minValue = value
							agents_execute = agentss
						}
					}
					ag.cout <- voteRequest{"execute", ag.name, ag.role, PingString, ag.cin, agents_execute, answer.cards[0:1], answer.Ja, game_vide, 0}
				} else {
					maxValue := math.MinInt32
					var agents_execute player
					for agentss, value := range ag.beliefs {
						if (value > maxValue) && (agentss.name != ag.name) {
							maxValue = value
							agents_execute = agentss
						}
					}
					ag.cout <- voteRequest{"execute", ag.name, ag.role, PingString, ag.cin, agents_execute, answer.cards[0:1], answer.Ja, game_vide, 0}
				}
			} else {
				fmt.Printf("Agent %q a reçu une demande de type %q. C'est incompréhensible.\n", ag.name, answer.typerequest)
			}

			// if answer.typerequest == "vote" {
			// 	fmt.Printf("reciu pongngngngn")
			// }

			time.Sleep(1 * time.Second)

		}
	}()
}
