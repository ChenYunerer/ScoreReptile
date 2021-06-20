package main

import (
	"ScoreReptile/src/job"
	"ScoreReptile/src/server"
)

func main() {
	job.StartCronJob()
	server.StartHttpServer()
}
