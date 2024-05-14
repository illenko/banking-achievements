package service

import (
	"fmt"
	"github.com/illenko/achievements-service/internal/mapper"
	"github.com/illenko/achievements-service/internal/model"
	"github.com/illenko/achievements-service/internal/repository"
	"github.com/illenko/achievements-service/pkg/http"
	"gofr.dev/pkg/gofr"
)

type AchievementService interface {
	GetAll(c *gofr.Context) ([]model.Achievement, error)
	GetResponse(c *gofr.Context) ([]http.Achievement, error)
	InsertAchievement(c *gofr.Context, a model.Achievement) error
	UpdateAchievement(c *gofr.Context, a model.Achievement) error
}

type achievementService struct {
	repo         repository.AchievementRepository
	settingsRepo repository.RuleRepository
	mapper       mapper.AchievementMapper
}

func NewAchievementService(repo repository.AchievementRepository, settingsRepo repository.RuleRepository, mapper mapper.AchievementMapper) AchievementService {
	return &achievementService{
		repo:         repo,
		settingsRepo: settingsRepo,
		mapper:       mapper,
	}
}

func (s *achievementService) GetAll(c *gofr.Context) ([]model.Achievement, error) {
	return s.repo.GetAll(c)
}

func (s *achievementService) GetResponse(c *gofr.Context) ([]http.Achievement, error) {
	achievements, err := s.repo.GetAll(c)
	if err != nil {
		return nil, err
	}

	allSettings, err := s.settingsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.mapper.ToResponse(achievements, allSettings)
}

func (s *achievementService) InsertAchievement(c *gofr.Context, a model.Achievement) error {
	err := s.repo.Insert(c, a)
	if err != nil {
		return fmt.Errorf("inserting achievement: %w", err)
	}
	return nil
}

func (s *achievementService) UpdateAchievement(c *gofr.Context, a model.Achievement) error {
	err := s.repo.Update(c, a)
	if err != nil {
		return fmt.Errorf("updating achievement: %w", err)
	}
	return nil
}
