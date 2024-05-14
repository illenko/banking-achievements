package mapper

import (
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/model"
	"github.com/illenko/achievements-service/pkg/http"
)

type AchievementMapper interface {
	ToResponse(achievements []model.Achievement, settings []model.Rule) ([]http.Achievement, error)
}

type achievementMapper struct{}

func NewAchievementMapper() AchievementMapper {
	return &achievementMapper{}
}

func (m *achievementMapper) ToResponse(achievements []model.Achievement, settings []model.Rule) ([]http.Achievement, error) {
	achievementMap := m.toAchievementMap(achievements)
	return m.toAchievementsResponse(settings, achievementMap), nil
}

func (m *achievementMapper) toAchievementMap(achievements []model.Achievement) map[uuid.UUID]model.Achievement {
	achievementMap := make(map[uuid.UUID]model.Achievement)
	for _, achievement := range achievements {
		achievementMap[achievement.SettingID] = achievement
	}
	return achievementMap
}

func (m *achievementMapper) toAchievementsResponse(settings []model.Rule, achievementMap map[uuid.UUID]model.Achievement) []http.Achievement {
	var responseAchievements []http.Achievement
	for _, setting := range settings {
		achievement, exists := achievementMap[setting.Id]
		if !exists {
			achievement = model.Achievement{
				ID:           uuid.New(),
				SettingID:    setting.Id,
				Value:        0,
				ValuesHolder: []string{},
			}
		}
		responseAchievements = append(responseAchievements, m.toAchievementResponse(setting, achievement))
	}
	return responseAchievements
}

func (m *achievementMapper) toAchievementResponse(setting model.Rule, achievement model.Achievement) http.Achievement {
	return http.Achievement{
		ID:          setting.Id,
		Name:        setting.Name,
		Description: setting.Description,
		Value:       int(achievement.Value),
		Goal:        setting.Criteria.Value,
		Repeatable:  setting.Repeatable,
		Count:       int(achievement.Value) / setting.Criteria.Value,
	}
}
