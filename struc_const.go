package main

const (
	Liberal = "Liberal"
	Fascist = "Fascist"
	Hitler  = "Hitler"
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
	deck                     []string //Les cartes du deck
	discard                  []string //Les cartes de la défausse
	logs                     []string
	liberalPolicies          int
	fascistPolicies          int
	investigationAvailable   bool //Si le pouvoir d'investigation est disponible
	specialElectionAvailable bool //Si le pouvoir d'élection spéciale est disponible
	executionAvailable       bool //Si le pouvoir d'exécution est disponible
	currentPresident         player //Président actuel
	currentChancellor        player //Chancelier actuel
	prevPresident            player //Président précédent
	prevChancellor           player //Chancelier précédent
	Hitler                   string
	website                  Website
	end                      bool
	propChancellor           player
	result_vote              bool
	nombre_echec             int //Nombre d'échecs de vote
	chaos                    bool
	result_game              string
	reponse                  string
}
