package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Shifts *gin.RouterGroup
}

// Register routes within gin router
func (r *Routes) Register(router *gin.Engine) {
	r.Shifts = router.Group("/shifts")
	{
		r.Shifts.GET("/", getAllShifts)
		r.Shifts.GET("/:id", getShiftByID)

		r.Shifts.POST("/", notImplemented)
		r.Shifts.POST("/:id", notImplemented)

		r.Shifts.PUT("/", notImplemented)
		r.Shifts.PUT("/:id", notImplemented)

		r.Shifts.PATCH("/", notImplemented)
		r.Shifts.PATCH("/:id", notImplemented)

		r.Shifts.DELETE("/", notImplemented)
		r.Shifts.DELETE("/:id", notImplemented)
	}
}

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}
