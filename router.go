package main

import (
	"context"
	"os"

	docs "github.com/cs301-2023-g3t3/points-ledger/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
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

	metadata := models.RequestMetadata{
		UserAgent: req.RequestContext.Identity.UserAgent,
		SourceIP:  req.RequestContext.Identity.SourceIP,
	}

	ctx = context.WithValue(ctx, "RequestMetadata", metadata)

	return ginLambda.ProxyWithContext(ctx, req)
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func InitRoutes() {
	router := gin.New()
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// router.Use(cors.Default())
	router.Use(middlewares.LoggingMiddleware())
  
	config := cors.DefaultConfig()
    config.AddAllowHeaders("Authorization", "X-IDTOKEN")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	health := new(controllers.HealthController)
	points := controllers.NewPointsController(*models.DB)

    docs.SwaggerInfo.BasePath = "docs"

	v1 := router.Group("/points")

	healthGroup := v1.Group("/health")
	healthGroup.GET("", health.CheckHealth)
	

	pointsGroup := v1.Group("/accounts")
	pointsGroup.GET("", points.GetAllAccounts)
	pointsGroup.GET("/paginate", points.GetPaginatedAccounts)
	pointsGroup.GET("/:ID", points.GetSpecificAccount)
	pointsGroup.GET("/user-account/:UserID", points.GetAccountByUser)

	pointsGroup.Use(middlewares.DecodeJWT())
	pointsGroup.PUT("/:ID", points.AdjustPoints)

    // Swagger
    router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	env := os.Getenv("ENV")
	if env == "lambda" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run()
	}
}
