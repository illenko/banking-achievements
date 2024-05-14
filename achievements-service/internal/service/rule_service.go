package service

import (
	"github.com/illenko/achievements-service/internal/model"
	"github.com/illenko/achievements-service/internal/repository"
)

type RuleService interface {
	GetAll() ([]model.Rule, error)
}

type ruleService struct {
	repo repository.RuleRepository
}

func NewRuleService(repo repository.RuleRepository) RuleService {
	return &ruleService{
		repo: repo,
	}
}

func (s *ruleService) GetAll() ([]model.Rule, error) {
	return s.repo.GetAll()
}
