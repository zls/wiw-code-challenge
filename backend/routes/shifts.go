package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getAllShifts(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func getShiftByID(c *gin.Context) {
	shiftID := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{})
}
