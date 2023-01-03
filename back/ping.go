package main

import (
	"fmt"
	"time"
)

const PingString = "ping"
const PongString = "pong"

type Request struct {
	typerequest string
	senderID    string
	req         string
	c           chan Request
}

type voteRequest struct {
	typerequest string
	senderID    string
	req         string
	c           chan string
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
	name  string
	role  string
	alive bool
	vote  string
	cin   chan string
	cout  chan voteRequest
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
	req.c <- PongString
}

func (g *game) handlePing(req voteRequest) {
	req.c <- PongString
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

func NewAgentPlayer(name string, cout chan voteRequest, role string, alive bool, vote string) *agentPlayer {
	cin := make(chan string)
	return &agentPlayer{name, role, alive, vote, cin, cout}
}

func (ag *agentPlayer) Start() {
	go func() {
		ag.cout <- voteRequest{"vote", ag.name, PingString, ag.cin}
		answer := <-ag.cin
		fmt.Printf("%q a pour role : %q\n", ag.name, ag.role)
		fmt.Printf("%q a reçu : %q\n", ag.name, answer)
		for {
			ag.cout <- voteRequest{"vote", ag.name, PingString, ag.cin}
			answer := <-ag.cin
			fmt.Printf("agent %q has received: %q\n", ag.name, answer)

			// if answer.typerequest == "vote" {
			// 	fmt.Printf("reciu pongngngngn")
			// }

			time.Sleep(10 * time.Second)

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
