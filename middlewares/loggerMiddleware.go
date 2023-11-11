package middlewares

import (
	// "encoding/json"
	// "fmt"
	"net/http"
	"time"

	"github.com/cs301-2023-g3t3/points-ledger/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CustomData struct {
	Method  string
	URL     string
	Headers map[string][]string
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()

		// Process the request
		ctx.Next()

		// End Time request
		endTime := time.Now()

		// Execution time
		latencyTime := endTime.Sub(startTime).Milliseconds()

		// Request data
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		// userAgent := ctx.GetHeader("User-Agent")
		metadata, ok := ctx.Request.Context().Value("RequestMetadata").(models.RequestMetadata)
		var userAgent string
		var sourceIP string
		if ok {
			// Access the UserAgent and SourceIP
			userAgent = metadata.UserAgent
			sourceIP = metadata.SourceIP
		}

		// Request IP

		if reqMethod == http.MethodPut {
			input, _ := ctx.Get("input")
			message, _ := ctx.Get("message")
			inputData := input.(models.Input)
			action := inputData.Action
			amount := inputData.Amount

			log.WithFields(log.Fields{
				"METHOD":     reqMethod,
				"URI":        reqUri,
				"STATUS":     statusCode,
				"LATENCY":    latencyTime,
				"USER_AGENT": userAgent,
				"SOURCE_IP":  sourceIP,
				"ACTION":     action,
				"AMOUNT":     amount,
				"MESSAGE":    message,
			}).Info("ADJUST POINTS REQUEST")
		}

		// if reqMethod == http.MethodGet {
		// 	log.WithFields(log.Fields{
		// 		"METHOD":     reqMethod,
		// 		"URI":        reqUri,
		// 		"STATUS":     statusCode,
		// 		"LATENCY":    latencyTime,
		// 		"USER_AGENT": userAgent,
		// 		"CLIENT_IP":  sourceIP,
		// 	}).Info("HTTP GET REQUEST")
		// }
	}
}
