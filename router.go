package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cs301-2023-g3t3/points-ledger/controllers"
	"github.com/cs301-2023-g3t3/points-ledger/middlewares"
	"github.com/cs301-2023-g3t3/points-ledger/models"
	"github.com/gin-contrib/cors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("agent: ", req.RequestContext.Identity.UserAgent)
	fmt.Println("ip add: ", req.RequestContext.Identity.SourceIP)
	metadata := models.RequestMetadata{
		UserAgent: req.RequestContext.Identity.UserAgent,
		SourceIP:  req.RequestContext.Identity.SourceIP,
	}

	// ctx = context.WithValue(ctx, "UserAgent", req.RequestContext.Identity.UserAgent)
	// ctx = context.WithValue(ctx, "SourceIP", req.RequestContext.Identity.SourceIP)
	ctx = context.WithValue(ctx, "RequestMetadata", metadata)

	return ginLambda.ProxyWithContext(ctx, req)
}

func InitRoutes() {
	router := gin.New()
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(middlewares.LoggingMiddleware())

	health := new(controllers.HealthController)
	points := controllers.NewPointsController(*models.DB)

	v1 := router.Group("/points")

	healthGroup := v1.Group("/health")
	healthGroup.GET("", health.CheckHealth)

	pointsGroup := v1.Group("/accounts")
	pointsGroup.GET("", points.GetAllAccounts)
	pointsGroup.GET("/:ID", points.GetSpecificAccount)
	pointsGroup.GET("/user-account/:UserID", points.GetAccountByUser)
	pointsGroup.PUT("/:ID", points.AdjustPoints)

	env := os.Getenv("ENV")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
