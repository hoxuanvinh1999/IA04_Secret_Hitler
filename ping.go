package main

import (
	"fmt"
	"time"
)

const PingString = "ping"
const PongString = "pong"

type Agent interface {
	Start()
}

type agentPlayer struct {
	name  string
	role  string
	alive bool
	vote  string
	cin   chan string
	cout  chan Request
}

type PongAgent struct {
	ID string
	c  chan Request
}

func (ag PongAgent) handlePing(req Request) {
	req.c <- PongString
}

func (ag *PongAgent) Start() {
	go func() {
		for {
			req := <-ag.c
			fmt.Printf("agent %q has received %q from %q\n", ag.ID,
				req.req, req.senderID)
			go ag.handlePing(req) // et si on enlève go ?
		}
	}()
}

func NewPongAgent(id string, c chan Request) *PongAgent {
	return &PongAgent{id, c}
}

func (ag *agentPlayer) Start() {
	go func() {
		for {
			ag.cout <- Request{ag.name, PingString, ag.cin}
			answer := <-ag.cin
			fmt.Printf("agent %q has received: %q\n", ag.name, answer)
			time.Sleep(time.Second)
		}
	}()
}

func NewAgentPlayer(name string, cout chan Request, role string, alive bool, vote string) *agentPlayer {
	cin := make(chan string)
	return &agentPlayer{name, role, alive, vote, cin, cout}
}

type Request struct {
	senderID string
	req      string
	c        chan string
}

func main() {
	c := make(chan Request)
	ponger := NewPongAgent("ponger", c)
	ponger.Start()
	for i := 0; i < 5; i++ {
		//id := fmt.Sprintf("pinger n°%d", i)
		name := []string{"Vinh", "Wassim", "Pierre", "Sylvain", "Nathan"}
		role := []string{Liberal, Fascist, Liberal, Hitler, Liberal}
		pinger := NewAgentPlayer(name[i], c, role[i], true, Liberal)
		pinger.Start()
	}
	time.Sleep(time.Minute)
}
