# IA04 - Secret Hitler
GHARBI Wassim - Nathan Le Boudec - Pierre Hannart - Xuan Vinh Ho

## Présentation du jeu 

Secret Hitler est un jeu dramatique d'intrigues et de trahisons politiques se déroulant dans
les années 1930 en Allemagne. Avant de lancer la partie, chaque joueur est aléatoirement et secrètement assigné à être un libéral ou un fasciste, et un joueur est Hitler. Les fascistes se coordonnent pour semer la méfiance et installer leur chef de sang froid ; les libéraux doivent trouver et arrêter Hitler avant qu'il ne soit trop tard.

AJOUTER UNE IMAGE DU PLATEAU DE JEU

### Rôles

5 à 6 joueurs : Les fascistes se connaissent entre eux. Mais seul Hitler sait qu'il est Hitler.

7 à 10 joueurs : Les fascistes, sauf Hitler, se connaissent entre eux. Hitler révèle son identité aux fascistes en faisant un signe de main les yeux fermés.

Les libéraux ne connaissent bien-sûr le parti d'aucun des joueurs.

### Conditions de victoire

Libéraux : 5 lois libérales promulguées OU Hitler exécuté.
Fascites : 6 lois fascistes promulguées OU Hitler élu Chancelier avec 3 lois fascistes promulguées.

### Règles spéciales

CHAOS : 
- La carte du dessus de la pioche LOIS est révélée et promulguée.
- Le jeton Election est réinitialisé.
- Un nouveau tour commence où tous les joueurs sont éligibles à devenir Chancelier.

Pouvoirs Présidentiels :
- Enquête : Le Président peut vérifier le parti d'un joueur. Il peut le partager aux autres s'il le souhaite.
- Election Spéciale : Interruption du tour. Le Président choisit le nouveau candidat à la Présidence. Celui-ci nomme un candidat à la Chancellerie. Le tour reprend.
- Aperçu des Lois : Le président regarde les 3 prochaines cartes de la poiche LOIS.
- Exécution : Le président tue un joueur.

### Tour de jeu

1) Le joueur à la gauche du Président devient le nouveau candidat à la Présidence.

2) Le candidat à la Présidence désigne un candidat à la Chancellerie.

3) Election : Les joueurs votent OUI ou NON pour le nouveau gouvernement. Il est élu si plus de la moitié votent OUI.

Si le NON gagne :

Le jeton Election n'est pas sur 3 --> Nouveau tour

Le jeton Election est sur 3 --> CHAOS

4) Tant que la loi n'est pas promulguée, le Président et le Chancelier doivent rester neutre.
- Le Président pioche 3 cartes LOIS et les regarde. Il en défausse une puis donne les 2 restantes au Chancelier.
- Le Chancelier regarde les cartes. Il en défausse une puis promulgue l'autre.
Cas particuliers :
- Si il y a déjà 5 lois fascites, le Chancelier peut utiliser son veto lorsqu'il reçoit les 2 cartes du Président. Si il est accepté par le Président, on commence un nouveau tour. Sinon la règle normale.
- La loi octroie un pouvoir. Le Président doit l'utiliser avant de passer au tour suivant.


## Objectifs du projet

L'objectif de notre projet est de simuler un partie de Secret-Hitler jouée par des agents, sans intervention externe.

Grâce à la simulation, nous tenterons de mettre en évidence la meilleur stratégie à adopter, selon les rôles, pour jouer à Secret-Hitler.

## Installation et Lancement du jeu

## Architecture et Code

## Discussions (differents points à noter sur le projet)

- Méta-paramètres 
- Intelligence des agents (comment ca impacte la simulation)
- Developper les situations que l'on peut rencontrer (victoire probable des fascistes ou liberaux, blocage,....)
- Niveau du code aussi

Le point influençant directement les résultats est le comportement et l'intelligence de jeu des agents. En effet, le jeu repose normalement beaucoup sur le débat entre les joueurs et la déduction. Par exemple, les fascites doivent mentir correctement sans incohérence tandis que les libéraux doivent déceler chaque indice. La manière de jouer, qui a été implémentée, reste donc basique.

A partir de ce point, nous avons donc pu faire ces observations : 




