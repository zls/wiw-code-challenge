package shifts

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	userID := c.Query("userID")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	shifts, err := model.GetShifts(c, userID, startTime, endTime)
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
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.Error(fmt.Errorf("failed to parse id to bytes, %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}
	idBytes, err := idUUID.MarshalBinary()
	if err != nil {
		c.Error(fmt.Errorf("failed to marshal uuid to bytes, %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	shift, err := model.GetShiftByID(c, idBytes)
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
	log.Printf("%+v", data)

	session := model.SessionFromContext(c)
	shift, err := model.NewShift(data.UserID, session.Account.ID, data.StartTime, data.EndTime)
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
		"shift": shift,
	})
}

func DeleteByID(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.Error(fmt.Errorf("failed to parse id to bytes, %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}
	idBytes, err := idUUID.MarshalBinary()
	if err != nil {
		c.Error(fmt.Errorf("failed to marshal uuid to bytes, %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}

	shift, err := model.GetShiftByID(c, idBytes)
	if err != nil {
		c.Error(fmt.Errorf("failed to get shift by id, %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{})
		return
	}
	err = model.DeleteShiftByID(c, idBytes, shift.UserID, shift.StartTime.Format("2006-01-02T15:04:05-0700"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}
