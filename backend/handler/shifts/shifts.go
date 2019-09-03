package shifts

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zls/wiw-code-challenge/backend/model"
)

func GetAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

type createForm struct {
	UID       int       `form:"uid"`
	AID       int       `form:"aid"`
	PID       int       `form:"pid"`
	StartTime time.Time `form:"startTime" time_format:"2006-01-02T15:04:05"`
	EndTime   time.Time `form:"endTime" time_format:"2006-01-02T15:04:05"`
}

func Create(c *gin.Context) {
	var data createForm
	err := c.BindJSON(&data)
	if err == nil {
		shift, err := model.NewShift(data.UID, data.PID, data.AID, data.StartTime, data.EndTime)
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
			"uid":       data.UID,
			"pid":       data.PID,
			"aid":       data.AID,
			"startTime": data.StartTime,
			"endTime":   data.EndTime,
		})
		return
	} else {
		log.Printf("error %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
}
