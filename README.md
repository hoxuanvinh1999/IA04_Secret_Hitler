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

### Qu'est-il sensé faire ?

L'objectif de notre projet est de simuler une partie de Secret Hitler jouée par des agents (plusieurs agents clients et un client serveur), sans intervention externe.
Par cette simulation de partie de jeu, nous cherchons à modéliser des comportements réalistes de joueurs en créant des agents faisant leurs propres choix en fonction des actions des autres joueurs et de leurs croyances.

### Que sommes nous sensés observer ?

Nous sommes sensés observer un plateau de jeu avec des joueurs (ayant des noms, des rôles, etc.) et effectuant des actions telles que poser des cartes.

### Problématique

Grâce à la simulation, nous tenterons de mettre en évidence la meilleur stratégie à adopter, selon les rôles, pour jouer (et gagner) à Secret Hitler.

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
Pour le back-end, nous avons 4 fichiers : main.go, player.go (agent client), gamemaster.go (agent serveur) et functions.go.
Le maître de jeu organise le bon déroulement de la partie et communique aux joueurs les résultats de leurs votes tandis que les joueurs prennent des décisions en fonction de leur rôle, leur stratégie et leurs croyances (qui évoluent compte tenu des comportements des autres joueurs). Il y a aussi une partie de hasard qui réside dans les attributions des rôles puis, tout du long de la partie, dans les cartes piochées par les joueurs et qui rend chaque partie absolument unique.
[à compléter]

Pour le front-end, Nous utilisons Gorilla WebSocket, qui est une implémentation Go du protocole WebSocket. Il créera un serveur et enverra des messages au site Web toutes les 200 ms. Le Web se met à jour automatiquement et lui montre comment le jeu fonctionne. 
Lien vers les [WebSocket](https://github.com/gorilla/websocket).

## Bilan

### Points positifs
1. Permet de suivre une partie de A à Z étape par étape et ainsi d'apprendre les règles du jeu très vite ou bien de décortiquer si l'on veut aller plus loin dans l'analyse de jeu.
2. La possibilité de faire varier des méta-paramètres tels que la stratégie adoptée par un joueur où la possibilité de connaître ou non les rôles des joueurs dès le début.
3. Permet de cerner les comportements typiques dans un jeu de déduction sociale.

### Points négatifs
1. Seule la version où les fascistes se connaissent entre eux et où seul Hitler sait qu'il est Hitler a été implémentée. Également, les règles spéciales n'ont pas été implémentées.
2. Le point influençant directement les résultats est le comportement est l'intelligence de jeu des agents. En effet, le jeu repose normalement beaucoup sur le débat entre les joueurs et la déduction. Par exemple, les fascistes doivent mentir correctement sans incohérence tandis que les libéraux doivent déceler chaque indice. La manière de jouer, qui a été implémentée, reste donc basique.
