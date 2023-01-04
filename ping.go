package main

import (
	"fmt"
	"time"
)

const PingString = "ping"
const PongString = "pong"

var c chan voteRequest
var player_vide = player{"rien", "rien", false, "rien", c}
var choice = ""

type Request struct {
	typerequest string
	senderID    string
	req         string
	c           chan Request
}

type voteRequest struct {
	typerequest string
	senderID    string
	role        string
	req         string
	c           chan voteRequest
	playerpres  player
	cards       []string
	Ja          bool
}

type Agent interface {
	Start()
}

type roleSend struct {
	role string
	c    chan string
}

type playerRequestToMj struct {
	playerToMj playerToMj
	senderID   string
	req        string
	c          chan string
}

type playerToMj struct {
	type_request       string //["prop_gov", "presidentdiscards", "chancellor_enact", "vote"]
	prop_gov           player
	president_discards string
	chancellor_enact   string
	vote               bool
}

type mjRequestToPlayer struct {
	mjToPlayer mjToPlayer
	senderID   string
	req        string
	c          chan string
}

type mjToPlayer struct {
	type_request       string //["prop_gov", "presidentdiscards", "chancellor_enact", "vote"]
	president_discards [3]string
	chancellor_enact   [2]string
	vote               player
}

type mjInformPlayer struct {
	liberalPolicies          int
	fascistPolicies          int
	investigationAvailable   bool
	specialElectionAvailable bool
	executionAvailable       bool
	currentPresident         string
	currentChancellor        string
}

type agentPlayer struct {
	name        string
	role        string
	alive       bool
	vote        string
	cin         chan voteRequest
	cout        chan voteRequest
	beliefs     map[player]int
	currentGame *game
}

type agentMJ struct {
	cin  chan Request
	cout chan Request
}

type PongAgent struct {
	ID string
	c  chan voteRequest
}

func (ag PongAgent) handlePing(req voteRequest) {
	req.c <- voteRequest{"choisischancelier", "MJ", PingString, "MJ", req.c, player_vide, []string{}, true}
}

func (ag PongAgent) choisisPres(req voteRequest) {
	req.c <- voteRequest{"choisischancelier", "MJ", PingString, "MJ", req.c, player_vide, []string{}, true}
}

func (g *game) choisisPres(req voteRequest) {
	req.c <- voteRequest{"choisischancelier", "MJ", PingString, "MJ", req.c, player_vide, []string{}, true}
}

func (g *game) handlePing(req voteRequest) {
	req.c <- voteRequest{"choisischancelier", "MJ", PingString, "MJ", req.c, player_vide, []string{}, true}
}

func NewPongAgent(id string, c chan voteRequest) *PongAgent {
	return &PongAgent{id, c}
}

func NewMJAgent(id string, c chan voteRequest) *PongAgent {
	return &PongAgent{id, c}
}

func (ag *PongAgent) Start() {
	go func() {
		for {
			req := <-ag.c
			fmt.Printf("agent %q has received %q from %q %q\n", ag.ID,
				req.req, req.senderID, req.typerequest)

			if req.typerequest == "vote" {
				fmt.Println("et c'est un vote")
			}
			go ag.handlePing(req) // et si on enlève go ?
		}
	}()
}

func NewAgentPlayer(name string, cout chan voteRequest, cin chan voteRequest, role string, alive bool, vote string, currentGame *game) *agentPlayer {
	//cin := make(chan voteRequest)
	beliefs := make(map[player]int, len(currentGame.players))
	for i := 0; i < len(currentGame.players); i++ {
		beliefs[currentGame.players[i]] = 3
	}

	return &agentPlayer{name, role, alive, vote, cin, cout, beliefs, currentGame}
}

func (ag *agentPlayer) Start(list_player []player) {
	go func() {
		for {
			//ag.cout <- voteRequest{"vote", ag.name, ag.role, PingString, ag.cin, player_vide}
			answer := <-ag.cin

			if answer.typerequest == "choisischancelier" {
				fmt.Printf("agent %q has received: %q\n", ag.name, answer.typerequest)
				fmt.Printf("Je vais choisir un chancelier.\n")
				for _, p := range list_player {
					if ag.role == Fascist || ag.role == Hitler {
						if (p.role == Fascist || p.role == Hitler) && (ag.name != p.name) {
							player_vide = p
							break
						}
					} else {
						if (ag.beliefs[p] > 2) && (ag.name != p.name) {
							player_vide = p
						}
					}
				}
				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, []string{}, true}

				//answer := <-ag.cin
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

				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, answer.cards, true}

			} else if answer.typerequest == "enact" {
				fmt.Printf("Agent %q a reçu une demande de %q avec pour cartes : %q\n", ag.name, answer.typerequest, answer.cards)
				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, answer.cards[0:1], true}

			} else if answer.typerequest == "vote" {
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

				ag.cout <- voteRequest{"prop_president", ag.name, ag.role, PingString, ag.cin, player_vide, answer.cards[0:1], answer.Ja}

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

// func main() {
// 	c := make(chan Request)
// 	ponger := NewPongAgent("ponger", c)
// 	ponger.Start()
// 	name := []string{"Vinh", "Wassim", "Pierre", "Sylvain", "Nathan"}
// 	role := []string{Liberal, Fascist, Liberal, Hitler, Liberal}
// 	for i := 0; i < 5; i++ {
// 		//id := fmt.Sprintf("pinger n°%d", i)
// 		pinger := NewAgentPlayer(name[i], c, role[i], true, Liberal)
// 		pinger.Start()
// 	}
// 	time.Sleep(time.Minute)
// }
