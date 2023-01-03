package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type Website struct {
	Game_title   string
	Time         string
	Player1_name string
	Player2_name string
	Player3_name string
	Player4_name string
	Player5_name string
	Player6_name string
}

func main() {

	go func() {
		website := Website{"Secret Hitler", time.Now().Format(time.Stamp), "Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"}
		template := template.Must(template.ParseFiles("web/main.html"))

		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if game_title := r.FormValue("white_background"); game_title != "" {
				website.Game_title = game_title
			}
			if err := template.ExecuteTemplate(w, "main.html", website); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
		fmt.Println(http.ListenAndServe(":8000", nil))
	}()

	go func() {
		c := make(chan voteRequest)
		g := newGame([]string{"Vinh", "Wassim", "Pierre", "Sylvain", "Jérôme", "Nathan"})
		g.c = c

		g.start()
		c_to_agent := make(map[string]chan voteRequest)
		//Création des agents joueurs
		for _, p := range g.players {
			newChan := make(chan voteRequest)
			c_to_agent[p.name] = newChan
			joueur := NewAgentPlayer(p.name, c, newChan, p.role, true, Liberal)
			joueur.Start(g.players)
		}
		g.c_to_agent = c_to_agent
		time.Sleep(10 * time.Minute)
	}()
	time.Sleep(10 * time.Minute)
}
