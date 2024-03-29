package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zls/wiw-code-challenge/backend/handler"
	"github.com/zls/wiw-code-challenge/backend/handler/shifts"
)

type Routes struct {
	Shifts *gin.RouterGroup
}

// Register routes within gin router
func (r *Routes) Register(router *gin.Engine) {
	r.Shifts = router.Group("/shifts")
	{
		r.Shifts.GET("/", shifts.GetShifts)
		r.Shifts.GET("/:id", shifts.GetByID)

		r.Shifts.POST("/", shifts.Create)

		r.Shifts.PUT("/:id", handler.NotImplemented)

		r.Shifts.DELETE("/:id", shifts.DeleteByID)
	}
}
