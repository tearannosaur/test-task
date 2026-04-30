package handlers

import (
	"app/service"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SubscriptionCreate godoc
// @Summary      Create subscription
// @Description  Create new subscription for user
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body service.SubscriptionRequest true "subscription data"
// @Success      201 {object} object{data=service.Subscription}
// @Failure      400 {object} object{error=string}
// @Failure      500 {object} object{error=string}
// @Router       /subscriptions [post]
func (h *Handler) SubscriptionCreate(c *gin.Context) {
	var req service.SubscriptionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		h.logger.Info("bad request", zap.Error(err))
		return
	}

	subscription, err := service.NewSubscription(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.logger.Error(err.Error(), zap.Error(err))
		return
	}

	err = h.repo.SubscriptionSave(subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create subscription",
		})
		h.logger.Error("failed to save subscription", zap.Error(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": subscription,
	})

}

// SubscriptionsList godoc
// @Summary      Get subscriptions list
// @Description  Returns list of all subscriptions
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Success      200 {object} object{data=[]service.SubscriptionResponse}
// @Failure      500 {object} object{error=string}
// @Router       /subscriptions [get]
func (h *Handler) SubscriptionsList(c *gin.Context) {
	list, err := h.repo.GetSubscriptionsList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		h.logger.Error("failed to get subscriptions list", zap.Error(err))
		return
	}
	responseList := service.ToResponseList(list)
	c.JSON(http.StatusOK, gin.H{
		"data": responseList,
	})
}

// GetSubscription godoc
// @Summary      Get subscription by id
// @Description  Returns single subscription by UUID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "subscription id (UUID)"
// @Success      200  {object}  object{data=service.SubscriptionResponse}
// @Failure      400  {object}  object{error=string}  "invalid uuid format"
// @Failure      404  {object}  object{error=string}  "id not found"
// @Failure      500  {object}  object{error=string}  "internal server error"
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	subId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid format",
		})
		h.logger.Info("invalid uuid format", zap.Error(err), zap.String("id", id))
		return
	}
	subscription, err := h.repo.GetSubscriptionById(subId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
			h.logger.Info("subscription not found", zap.String("id", subId.String()))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		h.logger.Error("failed to get subscription by id", zap.Error(err), zap.String("id", subId.String()))
		return
	}
	response := service.ToResponse(subscription)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// DeleteSubscription godoc
// @Summary      Delete subscription by id
// @Description  Deletes a subscription by UUID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "subscription id (UUID)"
// @Success      204  {string}  string  "no content"
// @Failure      400  {object}  object{error=string}  "invalid uuid format"
// @Failure      404  {object}  object{error=string}  "id not found"
// @Failure      500  {object}  object{error=string}  "internal server error"
// @Router       /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	subId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid format",
		})
		h.logger.Info("invalid uuid format", zap.Error(err), zap.String("id", id))
		return
	}
	err = h.repo.DeleteSubscriptionById(subId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
			h.logger.Info("subscription not found", zap.String("id", subId.String()))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		h.logger.Error("failed to delete subscription by id", zap.Error(err), zap.String("id", subId.String()))
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateSubscription godoc
// @Summary      Update subscription
// @Description  Update subscription by id
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "subscription id (UUID)"
// @Param        input  body      service.SubscriptionUpdateRequest  true  "updated subscription data"
// @Success      200    {object}  object{data=service.SubscriptionResponse}
// @Failure      400    {object}  object{error=string}  "invalid uuid / bad request / validation error"
// @Failure      404    {object}  object{error=string}  "id not found"
// @Failure      500    {object}  object{error=string}  "internal server error"
// @Router       /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	subId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid format",
		})
		h.logger.Info("invalid uuid format", zap.Error(err), zap.String("id", id))
		return
	}

	var req service.SubscriptionUpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		h.logger.Info("bad request", zap.Error(err))
		return
	}

	sub, err := service.UpdateSubscription(req, subId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.logger.Error(err.Error(), zap.Error(err))
		return
	}

	err = h.repo.UpdateSubscriptionById(sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
			h.logger.Info("subscription not found", zap.String("id", subId.String()))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		h.logger.Error("failed to update subscription by id", zap.Error(err), zap.String("id", subId.String()))
		return
	}

	subscription, err := h.repo.GetSubscriptionById(subId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
			h.logger.Info("subscription not found", zap.String("id", subId.String()))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		h.logger.Error("failed to get subscription by id", zap.Error(err), zap.String("id", subId.String()))
		return
	}
	response := service.ToResponse(subscription)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})

}

// GetTotal godoc
// @Summary      Get total subscription cost
// @Description  Calculate total cost of subscriptions for selected period with optional filters (user_id, service_name)
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        user_id      query     string  false  "user id (UUID)"
// @Param        service_name  query     string  false  "service name filter"
// @Param        from         query     string  true   "start period (MM-YYYY)"
// @Param        to           query     string  true   "end period (MM-YYYY)"
// @Success      200          {object}  object{total=int}  "total subscription cost"
// @Failure      400          {object}  object{error=string}  "invalid query params / bad date format / invalid uuid"
// @Failure      500          {object}  object{error=string}  "internal server error"
// @Router       /subscriptions/total [get]
func (h *Handler) GetTotal(c *gin.Context) {
	userIdStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, to, err := service.ParsePeriod(fromStr, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.logger.Error(err.Error(), zap.Error(err))
		return
	}
	var userId *uuid.UUID
	if userIdStr != "" {
		uid, err := uuid.Parse(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			h.logger.Info("invalid uuid format", zap.Error(err), zap.String("id", userIdStr))
			return
		}
		userId = &uid
	}
	subscription, err := h.repo.GetSubscriptionForPeriod(userId, serviceName, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		h.logger.Info("get subscriptions for period error", zap.Error(err))
		return
	}

	total := service.CountTotal(subscription, from, to)

	c.JSON(http.StatusOK, gin.H{
		"total": total,
	})
}
