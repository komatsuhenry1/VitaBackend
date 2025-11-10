package router

import (
	"medassist/internal/di"
	"medassist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(r *gin.RouterGroup, container *di.Container) {
	payment := r.Group("/payment")
	{
		payment.POST("/create-intent",middleware.AuthUser(), container.PaymentHandler.CreatePaymentIntent)
	}
}
