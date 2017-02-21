package ap

import (
	"bitbucket.org/dgub/evento/dom"
	"bitbucket.org/dgub/evento/usecases"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ReadEvent(c *gin.Context) {
	var err error

	streamid := c.Param("stream")
	if streamid == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid streamid"})
		return
	}

	var version uint64
	sVersion := c.Param("version")
	if sVersion == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid version"})
		return
	} else {
		version, err = strconv.ParseUint(sVersion, 10, 64)
		if err != nil {
			c.JSON(400, map[string]interface{}{"error": err.Error()})
			return
		}

	}

	event, err := usecases.ReadEvent(streamid, version)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}
	c.JSON(200, event)
}

func RangeQuery(c *gin.Context) {
	var err error

	streamid := c.Param("stream")
	if streamid == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid streamid"})
		return
	}

	var from, to uint64
	sFrom := c.Param("from")
	sTo := c.Param("to")
	if sFrom == "" || sTo == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid range"})
		return
	} else {
		from, err = strconv.ParseUint(sFrom, 10, 64)
		if err != nil {
			c.JSON(400, map[string]interface{}{"error": err.Error()})
			return
		}

		to, err = strconv.ParseUint(sTo, 10, 64)
		if err != nil {
			c.JSON(400, map[string]interface{}{"error": err.Error()})
			return
		}
	}

	list, err := usecases.RangeQuery(streamid, from, to)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"events": list})
}

func CreateEvent(c *gin.Context) {
	var e dom.Event
	err := c.BindJSON(&e)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}

	streamid := c.Param("stream")
	if streamid == "" {
		c.JSON(400, map[string]interface{}{"error": "invalid streamid"})
		return
	}
	e.StreamId = streamid

	var lock bool
	var version uint64
	eVersion := c.Request.URL.Query().Get("lock")
	if eVersion != "" {
		version, err = strconv.ParseUint(eVersion, 10, 64)
		if err != nil {
			c.JSON(400, map[string]interface{}{"error": err.Error()})
			return
		} else {
			lock = true
		}
	}

	var create bool
	created := c.Request.URL.Query().Get("create")
	if created == "true" {
		create = true
	}

	e.StreamId = streamid

	var id string
	if usecases.DbRunning() {
		id, version, err = usecases.EventoDb.StoreEvent(e, e.LinkStreams, lock, version, create)
		if err != nil {
			c.JSON(404, map[string]interface{}{"error": err.Error()})
			return
		}
	} else {
		c.JSON(302, map[string]interface{}{"error": "stoping db.."})
		return
	}

	c.JSON(200, map[string]interface{}{"id": id, "version": version})
}
