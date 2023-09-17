package main

import (
	"context"
	"os"

	"github.com/cs301-2023-g3t3/points-ledger/controllers"
	"github.com/cs301-2023-g3t3/points-ledger/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

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

	env := os.Getenv("env")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
