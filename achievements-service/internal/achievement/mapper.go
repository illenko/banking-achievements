package achievement

import (
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/settings"
	"github.com/illenko/achievements-service/pkg/http"
)

type Mapper interface {
	toResponse(a []Achievement, s []settings.Setting) ([]http.Achievement, error)
}

type mapper struct{}

func NewMapper() Mapper {
	return &mapper{}
}

func (m *mapper) toResponse(a []Achievement, s []settings.Setting) ([]http.Achievement, error) {

	achievementMap := make(map[uuid.UUID]Achievement)
	for _, achievement := range a {
		achievementMap[achievement.SettingID] = achievement
	}

	var achievements []http.Achievement
	for _, setting := range s {
		a, ok := achievementMap[setting.Id]
		if !ok {
			a = Achievement{
				ID:        uuid.New(),
				SettingID: setting.Id,
				Value:     0,
				Values:    []string{},
			}
		}

		achievements = append(achievements, http.Achievement{
			ID:          setting.Id,
			Name:        setting.Name,
			Description: setting.Description,
			Value:       int(a.Value),
			Goal:        setting.Criteria.Value,
			Repeatable:  setting.Repeatable,
			Count:       int(a.Value) / setting.Criteria.Value,
		})
	}

	return achievements, nil
}
