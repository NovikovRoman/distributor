package distributor

import (
	"strings"
)

var migrationMySQL = []string{
	"### 2021-09-27",
	"CREATE TABLE IF NOT EXISTS {table}",
	"(`event_id` VARCHAR(128) NOT NULL UNIQUE PRIMARY KEY,",
	"`message` VARCHAR(2048) NOT NULL,",
	"`created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP",
	") ENGINE = InnoDB;",
	/* Пример добавления миграций
	"### 2021-09-28",
	"IF NOT EXISTS(SELECT NULL FROM INFORMATION_SCHEMA.COLUMNS",
	"WHERE table_name = '{table}' AND table_schema = '{dbName}' AND column_name = 'newColumn')",
	"THEN ALTER TABLE `{table}` ADD `newColumn` BOOLEAN NOT NULL default 0;END IF;",*/
	"### end",
}

func (d MysqlDriver) Migration() (err error) {
	s := strings.Join(migrationMySQL, "\n")
	s = strings.Replace(s, "{table}", d.table, -1)
	s = strings.Replace(s, "{dbName}", d.dbName, -1)
	_, err = d.db.Queryx(s)
	return
}
