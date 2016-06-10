# coreingo

Pour le déploiement : 

Install go 1.6

`go get -u github.com/kataras/iris` & `go get -u github.com/FaXaq/gjp` (ce sont les seules dépendances du projet)
et ensuite `go get github.com/FaXaq/coreingo`
puis `go run *.go` pour lancer le fichier, sinon `go build *.go` pour compiler le fichier !
Installez les dépendances !
Voici les API disponibles :
* `POST /jobs` pour créer un job peut-être utilisé avec les arguments suivants :
  * `command` quelle commande à exécuter (`extract-audio` et `convert`)
  * `fromFile` chemin vers le fichier en input (ex: /home/user/Downloads/toto.mp4)
  * `toFile` nom et extension du fichier (ex: giphy.mkv)
  * Retourne un JSON contenant le job serialisé : 
`{
  "job": {
    "id": "9168fb1c-d488-4bdb-a0c8-ef97617bcf33",
    "name": "<path_to_file>/guitar to test.mp3",
    "status": "queued",
    "start": "0001-01-01T00:00:00Z",
    "end": "0001-01-01T00:00:00Z"
  }
}`
* `GET /jobs/search` donne une liste d'infos sur les jobs
  * `id` job id
  * Retourne un JSON contenant le job serialisé :
`{
   "id": "9168fb1c-d488-4bdb-a0c8-ef97617bcf33",
   "name": "<path_to_file>/guitar to test.mp3",
   "status": "proceeded",
   "start": "2016-03-04T12:37:20Z",
   "end": "0001-01-01T00:00:00Z"
}`
* `GET /jobs/progress` Arrête la jobPool
  * `id` job id
  * Retourne un JSON contenant uniquement la propriété percentage :
`{
   "percentage": 0.052427342648
}`

That's all folks 
