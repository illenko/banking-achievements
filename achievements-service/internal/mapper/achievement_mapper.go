package mapper

import (
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/model"
	"github.com/illenko/achievements-service/pkg/http"
)

type AchievementMapper interface {
	ToResponse(achievements []model.Achievement, rules []model.Rule) ([]http.Achievement, error)
}

type achievementMapper struct{}

func NewAchievementMapper() AchievementMapper {
	return &achievementMapper{}
}

func (m *achievementMapper) ToResponse(achievements []model.Achievement, rules []model.Rule) ([]http.Achievement, error) {
	achievementMap := m.toAchievementMap(achievements)
	return m.toAchievementsResponse(rules, achievementMap), nil
}

func (m *achievementMapper) toAchievementMap(achievements []model.Achievement) map[uuid.UUID]model.Achievement {
	achievementMap := make(map[uuid.UUID]model.Achievement)
	for _, achievement := range achievements {
		achievementMap[achievement.RuleID] = achievement
	}
	return achievementMap
}

func (m *achievementMapper) toAchievementsResponse(rules []model.Rule, achievementMap map[uuid.UUID]model.Achievement) []http.Achievement {
	var responseAchievements []http.Achievement
	for _, r := range rules {
		achievement, ok := achievementMap[r.ID]
		if !ok {
			achievement = model.Achievement{
				ID:           uuid.New(),
				RuleID:       r.ID,
				Value:        0,
				ValuesHolder: []string{},
			}
		}
		responseAchievements = append(responseAchievements, m.toAchievementResponse(r, achievement))
	}
	return responseAchievements
}

func (m *achievementMapper) toAchievementResponse(rule model.Rule, achievement model.Achievement) http.Achievement {
	return http.Achievement{
		ID:          rule.ID,
		Name:        rule.Name,
		Description: rule.Description,
		Value:       int(achievement.Value),
		Goal:        rule.Criteria.Value,
		Repeatable:  rule.Repeatable,
		Count:       int(achievement.Value) / rule.Criteria.Value,
	}
}
