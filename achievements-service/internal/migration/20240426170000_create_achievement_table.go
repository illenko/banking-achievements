package migration

import "gofr.dev/pkg/gofr/migration"

const createTable = `CREATE TABLE achievement (
    id UUID PRIMARY KEY,
    setting_id UUID NOT NULL,
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
