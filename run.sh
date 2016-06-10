export GJP_STARTJOB_CALLBACK=http://54.93.117.237:3000/notifications/task
export GJP_ENDJOB_CALLBACK=http://54.93.117.237:3000/notifications/task
export GJP_LOG_PATH=/var/log/core/
export GJP_WORK_PATH=/workdir/
export GJP_PORT=1337

go run *.go
