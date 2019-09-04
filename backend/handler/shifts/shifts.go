package shifts

import (
	"fmt"
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

func GetShifts(c *gin.Context) {
	shifts, err := model.GetShifts(c)
	if err != nil {
		c.Error(fmt.Errorf("failed to get shifts, %v", err))
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"shifts": shifts,
	})

}

func GetByID(c *gin.Context) {
	id := c.Param("id")
	shift, err := model.GetShiftByID(c, id)
	if err != nil {
		c.Error(fmt.Errorf("failed to get shift by id, %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"shift": shift,
	})
}

func Create(c *gin.Context) {
	var data model.Shift
	err := c.BindJSON(&data)
	if err != nil {
		c.Error(fmt.Errorf("unable to bind json to struct, %v", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	shift, err := model.NewShift(data.UserID, data.AccountID, data.StartTime, data.EndTime)
	if err != nil {
		c.Error(fmt.Errorf("unable to create new shift object %v", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}
	_, err = shift.Put(c)
	if err != nil {
		c.Error(fmt.Errorf("failed to write shift, %v", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shift": data,
	})
}
