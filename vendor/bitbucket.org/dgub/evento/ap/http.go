package ap

import (
	"github.com/gin-gonic/gin"
)

func InitHTTP(port string, dev bool) {
	if !dev {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.New()
	EventEndpoints(server)

	server.Run(":" + port)
}

func EventEndpoints(server *gin.Engine) {
	// Get stream current version
	server.GET("/stream/:stream", StreamVersion)

	// Read single event
	server.GET("/stream/:stream/:version", ReadEvent)

	// range query
	server.GET("/range/:stream/:from/:to", RangeQuery)

	// Push event into stream
	server.POST("/stream/:stream", CreateEvent)

	// Push event into stream
	server.DELETE("/stream/:stream", StreamDelete)

	// Push event into stream
	server.PUT("/stream/:stream/purge", StreamPurge)

	// Push event into stream
	server.GET("/purge/status", StreamToPurge)

	// Push event into stream
	server.PUT("/stream", StreamMatch)

	// create pub
	server.POST("/pub", CreatePub)

	// list pub
	server.GET("/pub", ListPub)

	// delete pub
	server.DELETE("/pub/:id", DeletePub)

	// server ping
	server.GET("/ping", Ping)
}
