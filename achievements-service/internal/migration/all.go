package migration

import "gofr.dev/pkg/gofr/migration"

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		20240426170000: createAchievementTable(),
	}
}
