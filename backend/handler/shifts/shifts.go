package shifts

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zls/wiw-code-challenge/backend/model"
)

func GetAll(c *gin.Context) {
	shifts, err := model.ScanShifts(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	c.JSON(http.StatusOK, gin.H{
		"shifts": shifts,
	})
}

func Create(c *gin.Context) {
	var data model.Shift
	err := c.BindJSON(&data)
	if err != nil {
		log.Printf("error %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	shift, err := model.NewShift(data.UserID, data.AccountID, data.StartTime, data.EndTime)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}
	_, err = shift.Put(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uid":       data.UserID,
		"aid":       data.AccountID,
		"startTime": data.StartTime,
		"endTime":   data.EndTime,
	})
}
