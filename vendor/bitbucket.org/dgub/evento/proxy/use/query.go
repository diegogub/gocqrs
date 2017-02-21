package use

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func QueryStream(c *gin.Context) {
	stream := c.Param("stream")
	if stream != "" {
		h, err := plog.AssignStream(stream)
		if err != nil {
			c.JSON(400, map[string]interface{}{"error": err.Error()})
		} else {
			c.JSON(200, h)
		}
		return
	}
}

func RegHost(c *gin.Context) {
	var host *EventoHost
	err := c.BindJSON(&host)
	if err != nil {
		return
	}

	err = host.Validate()
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}

	host = plog.AddHost(host)
	c.JSON(200, host)
}
