# coreingo

Basic implementation of ffmpeg CGI with webservices interface.

---

Deploy : 

Install go 1.6

`go get -u github.com/kataras/iris` & `go get -u github.com/FaXaq/gjp` (only dependencies)
and then `go get github.com/FaXaq/coreingo`

Launch the script `run.sh` to run.


----


API : 
* `POST /jobs` To create job :
  * `command` command to execute (`extract-audio` and `convert`)
  * `fromFile` path to file (ex: /home/user/Downloads/toto.mp4)
  * `toFile` name & ext for the output file (ex: giphy.mkv)
  * return : 
``` json
{
  "job": {
    "id": "9168fb1c-d488-4bdb-a0c8-ef97617bcf33",
    "name": "<path_to_file>/guitar to test.mp3",
    "status": "queued",
    "start": "0001-01-01T00:00:00Z",
    "end": "0001-01-01T00:00:00Z"
  }
}
```
* `GET /jobs/search` list of job info
  * `id` job id
  * return :
``` json
{
   "id": "9168fb1c-d488-4bdb-a0c8-ef97617bcf33",
   "name": "<path_to_file>/guitar to test.mp3",
   "status": "proceeded",
   "start": "2016-03-04T12:37:20Z",
   "end": "0001-01-01T00:00:00Z"
}
```
* `GET /jobs/progress` get job progress
  * `id` job id
  * return :
``` json
{
   "percentage": 0.052427342648
}
```

That's all folks 
