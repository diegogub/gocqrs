package ap

import (
	"bitbucket.org/dgub/evento/dom"
	"bitbucket.org/dgub/evento/usecases"
	"github.com/gin-gonic/gin"
)

func CreatePub(c *gin.Context) {
	var p dom.Pub

	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}

	err = usecases.CreatePub(p.Id, p.Regex, p.Url)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{"ok": true})
}

func ListPub(c *gin.Context) {
	list := usecases.ListPub()
	c.JSON(200, map[string]interface{}{"pubs": list})
}

func DeletePub(c *gin.Context) {
	id := c.Param("id")
	err := usecases.DeletePub(id)
	if err != nil {
		c.JSON(404, map[string]interface{}{"error": err.Error()})
	}

	c.JSON(200, map[string]interface{}{"ok": true})
}
