package shifts

import (
	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/zls/wiw-code-challenge/backend/model"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

type createForm struct {
	UID       int       `form:"uid"`
	AID       int       `form:"aid"`
	PID       int       `form:"pid"`
	StartTime time.Time `form:"startTime" time_format:"unix"`
	EndTime   time.Time `form:"endTime" time_format:"unix"`
}

func Create(c *gin.Context) {
	var data createForm
	if c.ShouldBindJSON(&data) == nil {
		c.JSON(http.StatusOK, gin.H{
			"uid":       data.UID,
			"pid":       data.PID,
			"aid":       data.AID,
			"startTime": data.StartTime,
			"endTime":   data.EndTime,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}
