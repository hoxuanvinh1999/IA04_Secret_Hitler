package main

import "sync"

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
	cin   chan voteRequest
}

type game struct {
	ID                       string
	c                        chan voteRequest
	c_to_agent               map[string]chan voteRequest
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
