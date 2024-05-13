package achievement

import (
	"fmt"
	"github.com/illenko/achievements-service/internal/settings"
	"github.com/illenko/achievements-service/pkg/http"
	"gofr.dev/pkg/gofr"
)

type Service interface {
	GetAll(c *gofr.Context) ([]http.Achievement, error)
}

type service struct {
	repo         Repository
	settingsRepo settings.Repository
	mapper       Mapper
}

func NewService(repo Repository, settingsRepo settings.Repository, mapper Mapper) Service {
	return &service{
		repo:         repo,
		settingsRepo: settingsRepo,
		mapper:       mapper,
	}
}

func (s *service) GetAll(c *gofr.Context) ([]http.Achievement, error) {
	achievements, err := s.repo.GetAll(c)
	if err != nil {
		return nil, err
	}

	allSettings, err := s.settingsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.mapper.toResponse(achievements, allSettings)
}

func (s *service) InsertAchievement(c *gofr.Context, a Achievement) error {
	err := s.repo.Insert(c, a)
	if err != nil {
		return fmt.Errorf("inserting achievement: %w", err)
	}
	return nil
}

func (s *service) UpdateAchievement(c *gofr.Context, a Achievement) error {
	err := s.repo.Update(c, a)
	if err != nil {
		return fmt.Errorf("updating achievement: %w", err)
	}
	return nil
}
