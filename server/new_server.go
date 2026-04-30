package server

import (
	"app/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "app/docs"

	"github.com/gin-gonic/gin"
)

func NewServer(h *handlers.Handler) error {
	r := gin.Default()
	r.POST("/subscriptions", h.SubscriptionCreate)
	r.GET("/subscriptions", h.SubscriptionsList)
	r.GET("/subscriptions/:id", h.GetSubscription)
	r.DELETE("/subscriptions/:id", h.DeleteSubscription)
	r.PUT("/subscriptions/:id", h.UpdateSubscription)
	r.GET("/subscriptions/total", h.GetTotal)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := r.Run()
	if err != nil {
		return err
	}
	return nil
}
