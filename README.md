# coreingo

Pour deploy : 

installer go !

`go get github.com/kataras/iris` & `go get github.com/FaXaq/iris`
et ensuite `go get github.com/FaXaq/coreingo` puis `go run go.go` pour lancer le fichier, sinon `go build go.go` pour compiler le fichier !
Installez les dépendances !
Et ... c'est à peu près tout, pour le moment ce que ça fait :

Le webserver est sur localhost:8000

* `/start` pour lancer la jobPool
* `/create` pour créer un job peut-être utilisé avec les arguments suivants :
  * `name` donner un nom à votre job !
  * `command` quelle commande sà exécuter
  * `args` les arguments ! un à la fois
  Exemple pour transcoder un fichier .mp4 en .mkv : 
  `http://localhost:8000/create?name=test&command=ffmpeg&args=-i&args=/home/marin/T%C3%A9l%C3%A9chargements/giphy.mp4&args=/home/marin/T%C3%A9l%C3%A9chargements/giphy.mkv`
  Retourne un id unique pour le job créé
* `/list` donne une liste d'infos sur les jobs
* `/search?id=iddujob` Donne des informations sur le job créé
* `/stop` Arrête la jobPool

TADAM
