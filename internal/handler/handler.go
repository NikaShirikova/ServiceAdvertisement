package handler

import (
	"advertisement/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
	log      *zap.Logger
	services *service.Service
	valid    *validator.Validate
	address  []string
	timeout  time.Duration
}

func NewHandler(s *service.Service, log *zap.Logger, valid *validator.Validate, adr []string) *Handler {
	return &Handler{
		services: s,
		log:      log.Named("handler"),
		valid:    valid,
		address:  adr,
		timeout:  time.Duration(time.Millisecond * 200),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/placements/request", h.PlacementRequest)

	return router
}
