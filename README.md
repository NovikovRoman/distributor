# Distributor

> Библиотека для распределения обработки событий между ботами

ID событий должны быть уникальными.

# Через MySQl или MariaDB

Запустить миграцию. Перед запуском бота и после обновлений библиотеки необходимо запускать миграцию
для установки возможных изменений в БД.

```go
package main

import (
	"github.com/NovikovRoman/distributor"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	var (
		db  *sqlx.DB
		err error
	)

	if db, err = sqlx.Connect("mysql", "<connectUrl>"); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	myDriver := distributor.NewMysqlDriver(db, "<dbName>", "<tableName>")
	if err = myDriver.Migration(); err != nil {
		log.Fatalln(err)
	}

	// …
}
```

Проверка должен ли бот обработать событие. Если `ok = false`, то `message` содержит в себе сообщение
от бота, который обработал событие.

```go
ok, message, err := myDriver.CanHandleEvent("<eventID>")
```

или

```go
ok, message, err := myDriver.CanHandleEvent(
"<eventID>", "<message1 for other bots>", "<message2 for other bots>")
```
