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
		r.Shifts.GET("/", shifts.GetAll)
		r.Shifts.GET("/:id", handler.NotImplemented)

		r.Shifts.POST("/", shifts.Create)
		r.Shifts.POST("/:id", handler.NotImplemented)

		r.Shifts.PUT("/", handler.NotImplemented)
		r.Shifts.PUT("/:id", handler.NotImplemented)

		r.Shifts.PATCH("/", handler.NotImplemented)
		r.Shifts.PATCH("/:id", handler.NotImplemented)

		r.Shifts.DELETE("/", handler.NotImplemented)
		r.Shifts.DELETE("/:id", handler.NotImplemented)
	}
}
