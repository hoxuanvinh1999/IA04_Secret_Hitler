# IA04 - Secret Hitler
Composition du groupe : Gharbi Wassim - Pierre Hannart - Xuan Vinh Ho - Nathan Le Boudec

## Présentation du jeu
Secret Hitler est un jeu à rôle secret (à l'instar du Loup Garou) dramatique d'intrigues et de trahisons politiques se déroulant dans les années 1930 en Allemagne. Avant de lancer la partie, chaque joueur est aléatoirement et secrètement assigné à être un libéral ou un fasciste, et un joueur est Hitler. Les fascistes se coordonnent pour semer la méfiance et installer leur chef de sang froid ; les libéraux doivent trouver et arrêter Hitler avant qu'il ne soit trop tard.

![](https://www.booksmith.com/sites/booksmith.com/files/styles/uc_product_full/public/SH_1.jpg?itok=EkNll-H7)
![](https://d1clhjx8k26u75.cloudfront.net/wp-content/uploads/2019/06/03170158/SecretHitlerBoardGameComponents.jpg)

### Rôles
5 à 6 joueurs : Les fascistes se connaissent entre eux. Mais seul Hitler sait qu'il est Hitler.

7 à 10 joueurs : Les fascistes, sauf Hitler, se connaissent entre eux. Hitler révèle son identité aux fascistes en faisant un signe de main les yeux fermés.

Les libéraux ne connaissent bien-sûr le parti d'aucun des joueurs.

### Conditions de victoire
Libéraux : 5 lois libérales promulguées OU Hitler exécuté.
Fascistes : 6 lois fascistes promulguées OU Hitler élu Chancelier avec 3 lois fascistes promulguées.

### Règles spéciales
Chaos : 
- La carte du dessus de la pioche LOIS est révélée et promulguée.
- Le jeton Election est réinitialisé.
- Un nouveau tour commence où tous les joueurs sont éligibles à devenir Chancelier.

Pouvoirs Présidentiels :
- Enquête : Le Président peut vérifier le parti d'un joueur. Il peut le partager aux autres s'il le souhaite.
- Election Spéciale : Interruption du tour. Le Président choisit le nouveau candidat à la Présidence. Celui-ci nomme un candidat à la Chancellerie. Le tour reprend.
- Aperçu des Lois : Le président regarde les 3 prochaines cartes de la poiche LOIS.
- Exécution : Le président tue un joueur.

### Tour de jeu
1) Le joueur à la gauche du Président devient le nouveau Président.

2) Le candidat à la Présidence désigne un candidat à la Chancellerie.

3) Election : Les joueurs votent OUI ou NON pour le nouveau gouvernement. Il est élu si plus de la moitié votent OUI.

Si le NON gagne :

Le jeton Election n'est pas sur 3 --> Nouveau tour

Le jeton Election est sur 3 --> CHAOS

4) Tant que la loi n'est pas promulguée, le Président et le Chancelier doivent rester neutres.
- Le Président pioche 3 cartes LOIS et les regarde. Il en défausse une puis donne les 2 restantes au Chancelier.
- Le Chancelier regarde les cartes. Il en défausse une puis promulgue l'autre.
Cas particuliers :
- Si il y a déjà 5 lois fascistes, le Chancelier peut utiliser son veto lorsqu'il reçoit les 2 cartes du Président. Si il est accepté par le Président, on commence un nouveau tour. Sinon la règle normale.
- La loi octroie un pouvoir. Le Président doit l'utiliser avant de passer au tour suivant.


## Objectifs du projet

### Qu'est-il sensé faire ?

L'objectif de notre projet est de simuler une partie de Secret Hitler jouée par des agents (plusieurs agents clients et un client serveur), sans intervention externe.
Par cette simulation de partie de jeu, nous cherchons à modéliser des comportements réalistes de joueurs en créant des agents faisant leurs propres choix en fonction des actions des autres joueurs et de leurs croyances.

### Que sommes nous sensés observer ?

Nous sommes sensés observer un plateau de jeu avec des joueurs (ayant des noms, des rôles, etc.) et effectuant des actions telles que poser des cartes.

### Problématique

Grâce à la simulation, nous tenterons de mettre en évidence quels traits de personnalité peuvent influer sur les résultats. Nous essaierons également de déterminer si le jeu est équilibré ou non.

## Installation et lancement du jeu
Pour exécuter le projet "Secret Hitler", vous devez suivre ces étapes :
1. Cloner ce projet dans votre machine propre: $ git clone https://gitlab.utc.fr/nleboude/ia04-secret-hitler --branch Branche_de_travail_Nathan
2. Entrer dans le dossier "ia04-secret-hitler" : $ cd ia04-secret-hitler/
Créer le module : $ go mod init ia04-secret-hitler
3. Ajouter le package go gorilla/websocket : $ go run .go get github.com/gorilla/websocket
4. Lancer les clients et le serveur : $ go run .  
5. Pour visualiser le front : ouvrez le fichier web/main.html dans votre dossier ia04-secret-hitler. Le jeu se lancera automatiquement, et vous pourrez voir les agents jouer.
6. Vous pouvez recommencer à partir de l'étape 4 pour relancer le jeu 

## Architecture
Pour le back-end, nous avons 4 fichiers : main.go, agent_player.go (agent client), game_master.go (agent serveur) et struc_cons.go.
Le maître de jeu organise le bon déroulement de la partie et communique aux joueurs les résultats de leurs votes tandis que les joueurs prennent des décisions en fonction de leur rôle, leur stratégie et leurs croyances (qui évoluent compte tenu des comportements des autres joueurs). Il y a aussi une partie de hasard qui réside dans les attributions des rôles puis, tout du long de la partie, dans les cartes piochées par les joueurs et qui rend chaque partie absolument unique.

Le fichier main.go lance une nouvelle partie avec un certain nombre de joueurs et leurs prénoms, fait en sorte que les joueurs fascistes se reconnaissent puis lance les agents joueurs (fonction Start du fichier agent_player.go).

Le fichier agent_player.go contient les fonctions permettant de créer et lancer les agents joueurs, de recevoir les informations transmises par le maître de jeu (cartes tirées, résultats d'élections, etc.) ainsi que transmettre des informations au maître de jeu (proposition d'un chancelier, votes, questions, etc.)

Le fichier game_master.go contient les fonctions permettant de créer et lancer le serveur, communiquer avec les agents clients (voir ci-dessus) ainsi que définir toutes les règles du jeu, gérer les piles de cartes ainsi que communiquer avec le front-end.

Le fichier struc_const.go contient des structures et constantes utilisées dans les autres fichiers.


Pour le front-end, nous utilisons Gorilla WebSocket, qui est une implémentation Go du protocole WebSocket. Il créera un serveur et enverra des messages au site Web toutes les 200 ms. Le Web se met à jour automatiquement et lui montre comment le jeu fonctionne. 
Lien vers les [WebSocket](https://github.com/gorilla/websocket).

Avec le serveur créé par Websocket, le reste est vraiment simple. Nous choisissons la voie simple en utilisant HTML et CSS pour le site Web. Nous utilisons ensuite Javascript pour traiter les messages qui sont envoyés par Websocket depuis Go. Le message est classé par la première lettre, Javascript traitera ces messages grâce à cela et changera le CSS. Avec cette méthode, le site Web sera modifié toutes les 200 ms car c'est l'écart entre 2 messages adjacents. Si nous l'envoyons plus rapidement ou plus lentement, certains changements du jeu peuvent être manqués.

## Détails des paramètres menteur et perspicacité

Comme demandé lors de la soutenance, voici une rapide explication de l'utilisation des paramètres menteur et perspicacité.
Lorsque l'agent candidat reçoit une question, il va tirer au sort un nombre, suivant une loi normale N(ag.menteur,1). En fonction de ce nombre, sa réponse sera plus ou moins suspecte.
Les autres agents vont ensuite être informés de cette réponse. Chaque agent va ensuite tirer au sort un nombre, suivant une loi normale N(ag.perspicacité,0.5). En fonction de cette valeur, il décèlera ou non si la réponse est suspecte, et à quel point elle l'est.
Les réalisations des variables aléatoires sont calculées à l'aide de cette fonction (merci RO05) :

```go
func RandomNormal(mean, stdDev float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return z0*stdDev + mean
}
```

## Bilan

Pour répondre à nos problématiques, nous avons choisi 9 types de parties différentes, avec différents profils de joueurs (bon, moyen et mauvais menteur (5, 2.5 et 0) et très perspicace, perspicace et peu perspicace (1, 0.5 et 0)). La statistique menteur est comprise entre 0 et 5, et la perspicacité entre 0 et 1. Nous avons réalisé 8 parties pour les 9 types de parties, soit 72 parties.

Voici les résultats :

|               | Bon menteur | Menteur moyen | Mauvais menteur |
|---------------|-------------|----------------|-----------------|
| Très perspicace | 50%        | 50%            | 62.5%           |
| Perspicace      | 37.5%      | 37.5%          | 50%             |
| Peu perspicace  | 12.5%      | 25%            | 25%             |


Nous arrivons aux conclusions que :
Si tout le monde ment mal, le taux de victoire des libéraux augmente et inversement.
Si tout le monde est perspicace, le taux de victoire des libéraux augmente et inversement.
Un agent fasciste mentant bien a plus de chance de gagner.
Un agent libéral perspicace a plus de chance de gagner.
Dans le cadre de parties avec de bons joueurs (bon menteurs et très perspicaces), le jeu est équilibré.
Cependant, avec des joueurs de niveaux variés, on a globalement plus de chance de gagner quand on est fasciste (61% de taux de réussite)
Quand le niveau baisse, le hasard a plus de place donc les fascistes ont encore plus de chance.

### Points positifs
1. Permet de suivre une partie de A à Z étape par étape et ainsi d'apprendre les règles du jeu très vite ou bien de décortiquer si l'on veut aller plus loin dans l'analyse de jeu.
2. Les agents sont cohérents dans leurs prises de décisions.
3. Utilisation de méta paramètres pour faire varier les profils des joueurs.
4. Permet de cerner les comportements typiques dans un jeu de déduction sociale.

### Points négatifs
1. Le raisonnement des agents pourrait être améliorée en utilisant notamment les probabilités d'apparition des cartes dans le deck.
2. Le point influençant directement les résultats est le comportement est l'intelligence de jeu des agents. En effet, le jeu repose normalement beaucoup sur le débat entre les joueurs et la déduction. Par exemple, les fascistes doivent mentir correctement sans incohérence tandis que les libéraux doivent déceler chaque indice. La manière de jouer, qui a été implémentée, reste donc basique.
3. Les joueurs communiquent seulement au moment du vote. Dans une partie, les échanges sont en permanence.
4. Dans ce genre de société, les affinités entre les agents peuvent entre en compte, il aurait été intéressant de les modéliser.
