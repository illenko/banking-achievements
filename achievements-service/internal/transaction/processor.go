package transaction

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/illenko/achievements-service/internal/achievement"
	"github.com/illenko/achievements-service/internal/settings"
	"gofr.dev/pkg/gofr"
	"slices"
)

type Processor interface {
	Process(c *gofr.Context) error
}

type processor struct {
	achievementService achievement.Service
	settingsService    settings.Service
}

func NewProcessor(achievementService achievement.Service, settingsService settings.Service) Processor {
	return &processor{
		achievementService: achievementService,
		settingsService:    settingsService,
	}
}

func (p *processor) Process(c *gofr.Context) error {
	var t Transaction
	err := c.Bind(&t)
	if err != nil {
		c.Logger.Error("Error unmarshalling transaction: ", err)
		return nil
	}
	c.Logger.Info("Consumed transaction", t)

	settings := filterAchievements(t, achievementSettings)
	achievements, err := fetchAchievementEntities(c)
	if err != nil {
		return err
	}

	return p.processAchievements(c, t, settings, achievements)
}

func (p *processor) processAchievements(c *gofr.Context, t Transaction, settings []settings.Setting, achievements []achievement.Achievement) error {
	achievementMap := make(map[uuid.UUID]achievement.Achievement)
	for _, a := range achievements {
		achievementMap[a.SettingID] = a
	}

	for _, setting := range settings {
		a, ok := achievementMap[setting.Id]
		if !ok {
			a = newAchievementEntity(setting.Id)
		}

		err := processAchievement(c, t, setting, a, ok)
		if err != nil {
			return fmt.Errorf("error processing achievement: %w", err)
		}
	}

	return nil
}

func newAchievementEntity(settingId uuid.UUID) achievement.Achievement {
	return achievement.Achievement{
		ID:        uuid.New(),
		SettingID: settingId,
		Value:     0,
		Values:    []string{},
	}
}

func (p *processor) processAchievement(c *gofr.Context, t Transaction, s settings.Setting, a achievement.Achievement, existing bool) error {
	switch s.Criteria.Type {
	case settings.Count:
		a.Value++
	case settings.Sum:
		a.Value += t.Amount
	case settings.Unique:
		value, err := getValue(s.Criteria, t)
		if err != nil {
			return fmt.Errorf("error getting value: %w", err)
		}
		if !slices.Contains(a.Values, value) {
			a.Values = append(a.Values, value)
			a.Value++
		}
	}

	if existing {
		err := p.achievementService.(c, a)
		if err != nil {
			return fmt.Errorf("error updating achievement entity: %w", err)
		}
	} else {
		err := createAchievementEntity(c, a)
		if err != nil {
			return fmt.Errorf("error creating achievement entity: %w", err)
		}
	}

	return nil
}

func getValue(criteria achievementCriteria, t transaction) (string, error) {
	if criteria.Field == Country {
		return t.Country, nil
	} else if criteria.Field == Category {
		return t.Category, nil
	} else {
		return "", fmt.Errorf("unknown field: %v", criteria.Field)
	}
}

func filterAchievements(t transaction, settings []achievementSetting) []achievementSetting {
	var filteredAchievements []achievementSetting

	for _, setting := range settings {
		if meetsCriteria(t, setting.Filter) {
			filteredAchievements = append(filteredAchievements, setting)
		}
	}

	return filteredAchievements
}

func meetsCriteria(t transaction, filter achievementFilter) bool {
	if filter.Amount > 0 && t.Amount < filter.Amount {
		return false
	}

	categoryMatched := false
	if filter.Categories != nil {
		for _, category := range *filter.Categories {
			if category == t.Category {
				categoryMatched = true
				break
			}
		}
	} else {
		categoryMatched = true
	}

	return categoryMatched
}
