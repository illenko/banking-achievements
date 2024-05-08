package migrations

import "gofr.dev/pkg/gofr/migration"

const createTable = `CREATE TABLE achievement_entity (
    id UUID PRIMARY KEY,
    value FLOAT NOT NULL,
    values TEXT[] NOT NULL
);`

func createAchievementTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTable)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
