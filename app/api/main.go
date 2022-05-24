// This is a simple http server (used Gin) to offer key APIs for approval frontend
// this service encompasses all apis need for a simplified transfer approval
// submit, reject, change, query
// this services use the workflow-client to interact with temporal
// as this is demo, we don't make this part very complicated only the necessary ones to demonstrate the features
package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO log and trace later add for demo
func main() {
	// Change Gin's Logging to a file.
	gin.DisableConsoleColor()
	// gin.DefaultWriter = io.MultiWriter(serverFile)
	r := gin.New()

	// middlewares
	r.Use(gin.Recovery())
	// r.Use(EnableTrace())
	// r.Use(AccessLog())

	// TODO find more usable format, add tracelog
	r.Use(gin.Logger())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// handlers bind to path
	RegisterRoutes(r)

	// TODO should change the port, so that not conflict with Temporal
	e := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if e != nil {
		panic(e)
	}
}
