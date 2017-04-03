package gocqrs

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type Viewer interface {
	ID() string
	Verb() string
	Generate(entity string, role, user string, c *gin.Context) (interface{}, error)
}
