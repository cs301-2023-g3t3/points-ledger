package main

import (
	"github.com/cs301-2023-g3t3/points-ledger/models"
)

func init() {
	models.ConnectToDB()
}

func main() {
	InitRoutes()
}
