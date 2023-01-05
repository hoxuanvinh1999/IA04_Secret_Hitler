package main

//ag1 suspecte fortement que ag2 est fasciste : 1
//ag1 suspecte que ag2 est fasciste : 2
//ag1 n'a pas de suspicions particulières concernant ag2 : 3
//ag1 suspecte que ag2 est libéral : 4
//ag1 suspecte fortement que ag2 est libéral : 5

func BeliefUp(ag1 *agentPlayer, ag2 player) {
	ag1.beliefs[ag2] += 1
}

func BeliefDown(ag1 *agentPlayer, ag2 player) {
	ag1.beliefs[ag2] -= 1
}

func (ag1 *agentPlayer) BothFascists(ag2 player) {
	for i := 0; i < len(ag1.currentGame.players); i++ {
		if ag1.currentGame.players[i] == ag2 {
			ag1.beliefs[ag1.currentGame.players[i]] = 1
		}
	}
}
