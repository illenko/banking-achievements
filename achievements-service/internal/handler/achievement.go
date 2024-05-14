package handler

import (
	"github.com/illenko/achievements-service/internal/service"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

type AchievementHandler interface {
	GetAll(c *gofr.Context) (interface{}, error)
}

type achievementHandler struct {
	service service.AchievementService
}

func NewAchievementHandler(service service.AchievementService) AchievementHandler {
	return &achievementHandler{
		service: service,
	}
}

func (h *achievementHandler) GetAll(c *gofr.Context) (interface{}, error) {
	res, err := h.service.GetResponse(c)

	if err != nil {
		return nil, err
	}

	return response.Raw{Data: res}, nil
}
