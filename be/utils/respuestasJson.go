package utils

import (
	"github.com/gin-gonic/gin"
)

func RespuestaJSON(c *gin.Context, status int, mensaje string, data ...interface{}) {
	resp := gin.H{
		"mensaje": mensaje,
	}

	if len(data) > 0 && data[0] != nil {
		resp["datos"] = data[0]
	}
	c.JSON(status, resp)
}
