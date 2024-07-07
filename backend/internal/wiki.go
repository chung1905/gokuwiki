package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListWiki(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"result": ListFiles(getPagesDir()),
	})
}
