package main

func main() {
	g := newGame([]string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"})

	// c := make(chan voteRequest)
	// ponger := NewPongAgent("ponger", c)

	// ponger.Start()
	// name := []string{"Vinh", "Wassim", "Pierre", "Sylvain", "Nathan"}
	// role := []string{Liberal, Fascist, Liberal, Hitler, Liberal}
	// for i := 0; i < 5; i++ {
	// 	//id := fmt.Sprintf("pinger n°%d", i)
	// 	pinger := NewAgentPlayer(name[i], c, role[i], true, Liberal)
	// 	pinger.Start()
	// }
	// time.Sleep(time.Second)

	//g.start(ponger)
	g.start()
}

//func (g *game) investigationAvailable() bool {
//	return g.fascistPolicies >= 3
//}
