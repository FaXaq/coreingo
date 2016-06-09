export GJP_STARTJOB_CALLBACK=http://localhost:3000/notifications/task
export GJP_ENDJOB_CALLBACK=http://localhost:3000/notifications/task
export GJP_LOG_PATH=.
export GJP_WORK_PATH=.
export GJP_PORT=1337

go run *.go
