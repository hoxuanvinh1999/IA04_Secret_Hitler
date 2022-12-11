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

type PongAgent struct {
	ID string
	c  chan Request
}

func (ag *PingAgent) Start() {
	go func() {
		for {
			ag.cout <- Request{ag.ID, PingString, ag.cin}
			answer := <-ag.cin
			fmt.Printf("agent %q has received: %q\n", ag.ID, answer)
			time.Sleep(20 * time.Second)
		}
	}()
}

func NewPingAgent(id string, cout chan Request) *PingAgent {
	cin := make(chan string)
	return &PingAgent{id, cin, cout}
}

type PingAgent struct {
	ID   string
	cin  chan string
	cout chan Request
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
		id := fmt.Sprintf("pinger n°%d", i)
		pinger := NewPingAgent(id, c)
		pinger.Start()
	}
	time.Sleep(time.Minute)
}
