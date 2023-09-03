package main

import (
	"github.com/cs301-2023-g3t3/points-ledger/controllers"
	"github.com/cs301-2023-g3t3/points-ledger/models"

	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)
	points := new(controllers.PointsController)

	v1 := router.Group("/api/v1")

	healthGroup := v1.Group("/health")
	healthGroup.GET("", health.CheckHealth)

	pointsGroup := v1.Group("/points")
	points.SetDB(models.DB)
	pointsGroup.GET("", points.GetAccounts)
	pointsGroup.GET("/:ID", points.GetSpecificAccount)
	pointsGroup.POST("/:ID", points.AdjustPoints)

	router.Run()
}	