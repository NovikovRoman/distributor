package distributor

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type mysqlEvent struct {
	EventID   string    `db:"event_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type MysqlDriver struct {
	db     *sqlx.DB
	dbName string
	table  string
}

//NewMysqlDriver driver for mysql.
func NewMysqlDriver(db *sqlx.DB, dbName, table string) *MysqlDriver {
	return &MysqlDriver{
		db:     db,
		dbName: dbName,
		table:  table,
	}
}

func (d MysqlDriver) CanHandleEvent(eventID string, message ...string) (ok bool, msg string, err error) {
	var (
		query  string
		set    string
		values string
	)

	evt := mysqlEvent{
		EventID:   eventID,
		Message:   strings.Join(message, " "),
		CreatedAt: time.Now(),
	}

	if set, values, err = fieldsForInsert(evt); err != nil {
		return
	}

	// случайная пауза
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Microsecond * time.Duration(rand.Int31n(1000000)))

	query = "INSERT INTO `" + d.table + "` (" + set + ") VALUES (" + values + ")"
	if _, err = d.db.NamedExec(query, evt); err == nil {
		ok = true
		msg = evt.Message

	} else if strings.Contains(err.Error(), "Duplicate entry") {
		query = "SELECT message FROM `" + d.table + "` WHERE `event_id` = ?"
		err = d.db.Get(&msg, query, eventID)
	}

	return
}

//fieldsForInsert возвращает поля для INSERT.
func fieldsForInsert(model interface{}) (set string, values string, err error) {
	var fields []string
	if fields, err = tableFields(model); err != nil {
		return
	}

	sqlValues := make([]string, len(fields))
	for i, name := range fields {
		fields[i] = "`" + name + "`"
		sqlValues[i] = ":" + name
	}

	set = strings.Join(fields, ",")
	values = strings.Join(sqlValues, ",")
	return
}

//tableFields возвращает имена полей таблицы.
func tableFields(values interface{}) (fields []string, err error) {
	v := reflect.ValueOf(values)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fields = []string{}
	switch {
	case v.Kind() == reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Tag.Get("db")
			if field == "-" {
				continue

			} else if field == "" {
				fields = append(fields, strings.ToLower(v.Type().Field(i).Name))
				continue
			}

			fields = append(fields, field)
		}
		return

	case v.Kind() == reflect.Map:
		fields = make([]string, len(v.MapKeys()))
		for i, k := range v.MapKeys() {
			fields[i] = k.String()
		}
		return
	}

	err = fmt.Errorf("dbTableFields requires a struct or a map, found: %s", v.Kind().String())
	return
}
