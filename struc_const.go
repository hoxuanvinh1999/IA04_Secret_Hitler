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

type Website struct {
	Game_title         string
	Time               string
	Players_name       []string
	Players_side       []string
	Players_alive      []string
	Liberal_board      []string
	Fascist_board      []string
	Current_President  string
	Current_Chancellor string
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
	currentPresident         player
	currentChancellor        player
	prevPresident            player
	prevChancellor           player
	Hitler                   string
	website                  Website
	end                      bool
}
