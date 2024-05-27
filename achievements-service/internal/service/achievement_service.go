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
	InsertAchievement(c *gofr.Context, achievement model.Achievement) error
	UpdateAchievement(c *gofr.Context, achievement model.Achievement) error
}

type achievementService struct {
	repo        repository.AchievementRepository
	ruleService RuleService
	mapper      mapper.AchievementMapper
}

func NewAchievementService(repo repository.AchievementRepository, ruleService RuleService, mapper mapper.AchievementMapper) AchievementService {
	return &achievementService{
		repo:        repo,
		ruleService: ruleService,
		mapper:      mapper,
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

	rules, err := s.ruleService.GetAll()
	if err != nil {
		return nil, err
	}

	return s.mapper.ToResponse(achievements, rules)
}

func (s *achievementService) InsertAchievement(c *gofr.Context, achievement model.Achievement) error {
	err := s.repo.Insert(c, achievement)
	if err != nil {
		return fmt.Errorf("inserting achievement: %w", err)
	}
	return nil
}

func (s *achievementService) UpdateAchievement(c *gofr.Context, achievement model.Achievement) error {
	err := s.repo.Update(c, achievement)
	if err != nil {
		return fmt.Errorf("updating achievement: %w", err)
	}
	return nil
}
