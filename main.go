package main

import (
	"github.com/cs301-2023-g3t3/points-ledger/models"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	models.ConnectToDB()
}

func main() {
	InitRoutes()
}
