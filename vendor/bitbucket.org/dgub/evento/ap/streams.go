package ap

import (
	/*
		"bitbucket.org/dgub/apis"
		"bitbucket.org/dgub/evento/dom"
	*/
	"bitbucket.org/dgub/evento/usecases"
	"github.com/gin-gonic/gin"
)

func StreamVersion(c *gin.Context) {

	streamid := c.Param("stream")
	if streamid == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid streamid"})
		return
	}

	version, err := usecases.StreamVersion(streamid)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]interface{}{"version": version})
}

type StreamQuery struct {
	Regex string `json:"regex"`
}

func StreamMatch(c *gin.Context) {
	var sq StreamQuery
	err := c.BindJSON(&sq)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
		return
	}
	list, err := usecases.StreamMatch(sq.Regex)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"streams": list})
}

func StreamDelete(c *gin.Context) {
	streamid := c.Param("stream")
	if streamid == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid streamid"})
		return
	}

	err := usecases.EventoDb.DeleteStream(streamid)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"stream": streamid, "deleted": true})
}

// Purge deleted stream
func StreamPurge(c *gin.Context) {
	streamid := c.Param("stream")
	if streamid == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid streamid"})
		return
	}

	err := usecases.EventoDb.PurgeStream(streamid)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"stream": streamid, "deleted": true})
}

func StreamToPurge(c *gin.Context) {
	list := usecases.GetPurges()
	c.JSON(200, map[string]interface{}{"streams": list})
}
