# coreingo

Pour deploy : 

installer go 1.6 mini

`go get -u github.com/kataras/iris` & `go get -u github.com/FaXaq/gjp` et ensuite `go get github.com/FaXaq/coreingo` puis `go run *.go` pour lancer le fichier, sinon `go build *.go` pour compiler le fichier !
Installez les dépendances !
Et ... c'est à peu près tout, pour le moment ce que ça fait :

* `POST /jobs` pour créer un job peut-être utilisé avec les arguments suivants :
  * `command` quelle commande sà exécuter (`extract-audio` et `convert`)
  * `fromFile` chemin vers le fichier en input (ex: /home/user/Downloads/toto.mp4)
  * `toFile` nom et extension du fichier (ex: giphy.mkv)
  Retourne un id unique pour le job créé
* `GET /jobs/search` donne une liste d'infos sur les jobs
  * `id` job id
* `GET /jobs/progress` Arrête la jobPool
  * `id` job id


TADAM
